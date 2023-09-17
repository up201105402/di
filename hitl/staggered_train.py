# Imports
import os
import numpy as np
from PIL import ImageFile
from tqdm import tqdm

# Sklearn and SciPy Imports
from sklearn.metrics import accuracy_score, recall_score, precision_score, f1_score
from scipy.stats import entropy

# PyTorch Imports
import torch
import torchvision.transforms as transforms

# Project Imports
from xai_utilities import takeThird, GenerateDeepLiftAtts
from ui_utilities import GetOracleFeedback, matchSelectedRects

# Global variables and definitions
ImageFile.LOAD_TRUNCATED_IMAGES = True
HITL_LAMBDA = 1e7

# Function: Custom loss function
def my_loss(Ypred, X, W):
    # it "works" with both retain_graph=True and create_graph=True, but I think
    # the latter is what we want
    grad = torch.autograd.grad(Ypred.sum(), X, retain_graph=True)[0]
    # grad has a shape: torch.Size([256, 3, 32, 32])
    
    return (W *(grad**2).mean(1)).mean()

# Train model and sample the most useful images for decision making (entropy based sampling)
def staggered_active_train_model(
    model, 
    model_name, 
    train_loader, 
    val_loader, 
    history_dir, 
    weights_dir, 
    epochs_dir, 
    entropy_thresh, 
    nr_queries, 
    start_epoch, 
    data_classes, 
    oversample, 
    sampling_process, 
    EPOCHS, 
    DEVICE, 
    LOSS, 
    percentage=100, 
    resume_epoch=0, 
    should_resume=False, 
    image_resize_factor=1,
    optimizer=torch.optim.Adam,
    learning_rate=1e-5
    ):
    
    assert sampling_process in ['low_entropy', 'high_entropy']
    # Hyper-parameters
    LEARNING_RATE=learning_rate
    OPTIMISER = optimizer(model.parameters(), lr=LEARNING_RATE)

    # Initialise min_train and min_val loss trackers
    min_train_loss = np.inf
    min_val_loss = np.inf

    # Initialise losses arrays
    train_losses = np.zeros((EPOCHS, ))
    val_losses = np.zeros_like(train_losses)

    # Initialise metrics arrays
    train_metrics = np.zeros((EPOCHS, 4))
    val_metrics = np.zeros_like(train_metrics)

    # Weights for human in the loop loss
    print(f"Length of train loader: {len(train_loader)}")
    W = torch.zeros((len(train_loader.dataset), 224, 224), device=DEVICE)

    for epoch in range(resume_epoch, EPOCHS, 1):
        # Epoch 
        print(f"Epoch: {epoch+1}")
        
        # Training Loop
        print(f"Training Phase")
        
        # Initialise lists to compute scores
        y_train_true = list()
        y_train_pred = list()

        # Initialise list of high entropy predictions
        # and corresponding images
        informative_pred = list()

        # Running train loss
        run_train_loss = 0.0
        vanilla_run_train_loss = 0.0

        epoch_dir = os.path.join(epochs_dir, str(epoch))

        # Put model in training mode
        model.train()

        if (resume_epoch == 0 and not should_resume) or epoch != resume_epoch:
            # Iterate through dataloader
            for batch_idx, (images, images_og, labels, indices) in enumerate(tqdm(train_loader)):
                # move data, labels and model to DEVICE (GPU or CPU)
                if W.sum() > 0 and oversample == True:
                    # print('W sum 1,2:', W.sum([1, 2]) > 0)
                    # print('len(W):', len(W))
                    ix_annotated = torch.arange(0, len(W)).to(DEVICE)[W.sum([1, 2]) > 0]
                    # ix_annotated = [0, 4, 9, ...]
                    # print('ix_annotated:', ix_annotated)
                    i = ix_annotated[torch.randint(0, len(ix_annotated), ())]
                    # print('adicionei a amostra', i)
                    _image, _image_og, _label, _indice = train_loader.dataset[i]
                    images = torch.cat((images, _image[None]))
                    images_og = torch.cat((images_og, _image_og[None]))
                    labels = torch.cat((labels, torch.tensor([_label])))
                    indices = torch.cat((indices, torch.tensor([_indice])))

                images, labels, images_og = images.to(DEVICE), labels.to(DEVICE), images_og.to(DEVICE)
                model = model.to(DEVICE)

                # Find the loss and update the model parameters accordingly
                # Clear the gradients of all optimized variables
                OPTIMISER.zero_grad()

                # Forward pass: compute predicted outputs by passing inputs to the model
                images.requires_grad = True
                logits = model(images)
                loss = LOSS(logits, labels)

                vanilla_loss = loss.item()
                images_og.requires_grad = True
                logits = model(images_og)
                custom_loss = my_loss(logits, images_og, W[indices])*HITL_LAMBDA
                loss += custom_loss
                
                
                #if (epoch >= start_epoch):
                #    print(f"Cross entropy loss: {vanilla_run_train_loss}")
                #    print(f"Loss After AL: {loss} ")
                #    print(f"Custom imposed loss: {custom_loss}")

                # Backward pass: compute gradient of the loss with respect to model parameters
                loss.backward()
                
                # Perform a single optimization step (parameter update)
                OPTIMISER.step()
                
                # Initialise Active Learning
                if(epoch >= start_epoch and epoch < 20):
                    # Copy logits to cpu

                    #with torch.no_grad():
                    #    _logits = model(images_og)

                    pred_logits = logits.cpu().detach().numpy()
                    pred_logits = torch.FloatTensor(pred_logits)
                    pred_probs = torch.softmax(pred_logits,1)
                    _,pred = torch.max(pred_probs, 1)
                    #print('predictions:', pred_probs.shape, pred)
                    
                    if(sampling_process == 'high_entropy'):
                        # Iterate logits tensor 
                        for idx in range(len(pred_probs)):
                            # calculate entropy for each single image logits in batch
                            pred_entropy = entropy(pred_probs[idx])
                            if(pred_entropy > entropy_thresh):
                                temp_image_info = [images_og[idx], labels[idx], pred_entropy, indices[idx], idx, pred[idx]]
                                informative_pred.append(temp_image_info)
                    elif(sampling_process == 'low_entropy'):
                        for idx in range(len(pred_probs)):
                            pred_entropy = entropy(pred_probs[idx])
                            
                            if pred_entropy < entropy_thresh and pred[idx] != labels[idx].cpu():
                                temp_image_info = [images_og[idx], labels[idx], pred_entropy, indices[idx], idx, pred[idx]]
                                informative_pred.append(temp_image_info)
                                
                        
                
                # Update batch losses
                run_train_loss += (loss.item() * images.size(0))
                vanilla_run_train_loss += (vanilla_loss*images.size(0))

                # Concatenate lists
                y_train_true += list(labels.cpu().detach().numpy())
                
                # Using Softmax
                # Apply Softmax on Logits and get the argmax to get the predicted labels
                s_logits = torch.nn.Softmax(dim=1)(logits)
                s_logits = torch.argmax(s_logits, dim=1)
                y_train_pred += list(s_logits.cpu().detach().numpy())

            # Compute Average Train Loss
            avg_train_loss = run_train_loss/len(train_loader.dataset)
            vanilla_avg_train_loss = vanilla_run_train_loss/len(train_loader.dataset)

            # Compute Train Metrics
            train_acc = accuracy_score(y_true=y_train_true, y_pred=y_train_pred)
            train_recall = recall_score(y_true=y_train_true, y_pred=y_train_pred, average="weighted")
            train_precision = precision_score(y_true=y_train_true, y_pred=y_train_pred, average="weighted")
            train_f1 = f1_score(y_true=y_train_true, y_pred=y_train_pred, average="weighted")

            # Get highest entropy prediction information
            if(sampling_process == 'high_entropy' and len(informative_pred) > 0):
                informative_pred.sort(key=takeThird, reverse=True)
                print(f"Highest entropy predictions after {epoch+1} epochs: ")
            elif(sampling_process == 'low_entropy' and len(informative_pred) > 0):
                informative_pred.sort(key=takeThird, reverse=False)
                print(f"Lowest entropy predictions after {epoch+1} epochs: ")
            
            if len(informative_pred) > 0:
                if not os.path.isdir(epoch_dir):
                    os.makedirs(epoch_dir)

            # Print query entropies and perform Deep Lift on each data point

            for i in range(len(informative_pred)):
                if(i < nr_queries):
                    print(informative_pred[i][2]) 
                    query_image = informative_pred[i][0]
                    query_index = informative_pred[i][4]
                    image_index = informative_pred[i][3]
                    query_label = informative_pred[i][1]
                    temp_pred = informative_pred[i][5]
                    deepLiftAtts, query_pred = GenerateDeepLiftAtts(image=query_image, label=query_label, model = model, data_classes=data_classes, temp_pred=temp_pred)
                    
                    print("query_pred: ",query_pred,"query_label", query_label)
                    # Aggregate along color channels and normalize to [-1, 1]
                    deepLiftAtts = deepLiftAtts.sum(axis=np.argmax(np.asarray(deepLiftAtts.shape) == 3))
                    deepLiftAtts /= np.max(np.abs(deepLiftAtts))
                    deepLiftAtts = torch.tensor(deepLiftAtts)
                    #print(deepLiftAtts.shape)

                    __, ___ = GetOracleFeedback(image=query_image.detach().cpu().numpy(), label=query_label, idx=image_index, model_attributions=deepLiftAtts, pred=query_pred, rectSize=28, rectStride=28, nr_rects=10, epoch_number=epoch, query_nr=i, image_resize_factor=image_resize_factor, epoch_dir=epoch_dir)
                    # print(selectedRectangles)

                    # change the weights W=1 in the selected rectangles area
                    # print("index:", image_index)
                    # print(f"Length of rectangle vector: {len(W)}")
                    # for rect in selectedRectangles:
                    #     W[image_index, rect[1]:rect[3], rect[0]:rect[2]] = 1
        
        image_indexes = [t[3].item() for t in informative_pred]

        if len(informative_pred) > 0 and nr_queries > 0:
            # Save needed information
            np.save(f"{epoch_dir}/image_indexes.npy", np.asarray(image_indexes))
            torch.save(W, f"{epoch_dir}/W.pt")
            np.save(f"{epoch_dir}/min_train_loss.npy", min_train_loss)
            np.save(f"{epoch_dir}/min_val_loss.npy", min_val_loss)
            np.save(f"{epoch_dir}/train_losses.npy", train_losses)
            np.save(f"{epoch_dir}/val_losses.npy", val_losses)
            np.save(f"{epoch_dir}/train_metrics.npy", train_metrics)
            np.save(f"{epoch_dir}/val_metrics.npy", val_metrics)
            np.save(f"{epoch_dir}/vanilla_avg_train_loss.npy", vanilla_avg_train_loss)
            np.save(f"{epoch_dir}/train_acc.npy", train_acc)
            np.save(f"{epoch_dir}/train_recall.npy", train_recall)
            np.save(f"{epoch_dir}/train_precision.npy", train_precision)
            np.save(f"{epoch_dir}/train_f1.npy", train_f1)

            # Save checkpoint
            model_path = os.path.join(epoch_dir, f"{model_name}_{percentage}p_{EPOCHS}e_{sampling_process}_epoch_{epoch}.pt")
            torch.save(model.state_dict(), model_path)

            print(f"Successfully saved at: {model_path}")
            return val_losses,train_losses,val_metrics,train_metrics
        
        if resume_epoch == epoch and should_resume:
            # Load needed information
            image_indexes = np.load(f"{epoch_dir}/image_indexes.npy")
            W = torch.load(f"{epoch_dir}/W.pt")
            min_train_loss = np.load(f"{epoch_dir}/min_train_loss.npy")
            min_val_loss = np.load(f"{epoch_dir}/min_val_loss.npy")
            train_losses = np.load(f"{epoch_dir}/train_losses.npy")
            val_losses = np.load(f"{epoch_dir}/val_losses.npy")
            train_metrics = np.load(f"{epoch_dir}/train_metrics.npy")
            val_metrics = np.load(f"{epoch_dir}/val_metrics.npy")
            vanilla_avg_train_loss = np.load(f"{epoch_dir}/vanilla_avg_train_loss.npy")
            train_acc = np.load(f"{epoch_dir}/train_acc.npy")
            train_recall = np.load(f"{epoch_dir}/train_recall.npy")
            train_precision = np.load(f"{epoch_dir}/train_precision.npy")
            train_f1 = np.load(f"{epoch_dir}/train_f1.npy")
        
        selectedRectangles = {}

        # load selected rectangles
        for i in range(len(image_indexes)):
            if(i < nr_queries):
                try:
                    selectedRects = np.genfromtxt(f"{epoch_dir}/query_{i}_rects_selected.csv", dtype=np.int64, delimiter=",", ndmin=2)
                    rectsTensors = torch.load(f"{epoch_dir}/rects_{i}.pt")
                    selectedRectangles = matchSelectedRects(selectedRects, rectsTensors)

                    image_index = image_indexes[i]

                    # change the weights W=1 in the selected rectangles area
                    print("index:", image_index)
                    print(f"Length of rectangle vector: {len(W)}")
                    for rect in selectedRectangles:
                        W[image_index, rect[1]:rect[3], rect[0]:rect[2]] = 1
                except FileNotFoundError:
                    continue

        # Print high entropy prediciton data points
        print(f"Number of informative predictions after {epoch+1} epochs: {len(image_indexes)}")

        # # Visualize entropy distribution
        # df = pd.DataFrame(informative_pred, columns = ['Image','Label','Entropy'])
        # plt.hist(df['Entropy'], color = 'blue', edgecolor = 'black',
        #  bins = int(160/10))
        # plt.savefig(f"/home/up201605633/Desktop/Results/DeepLift/AL_tests/entropy_dist_e{epoch+1}.png")


        
        # Print Statistics
        print(f"Train Loss: {vanilla_avg_train_loss}\tTrain Accuracy: {train_acc}")
        
        # Append values to the arrays
        # Train Loss
        train_losses[epoch] = vanilla_avg_train_loss
        
        # Save it to directory
        fname = os.path.join(history_dir, f"{model_name}_train_losses_{percentage}_{sampling_process}.npy")
        np.save(file=fname, arr=val_losses, allow_pickle=True)

        
        # Train Metrics
        # Acc
        train_metrics[epoch, 0] = train_acc
        # Recall
        train_metrics[epoch, 1] = train_recall
        # Precision
        train_metrics[epoch, 2] = train_precision
        # F1-Score
        train_metrics[epoch, 3] = train_f1
        
        fname = os.path.join(history_dir, f"{model_name}_train_metrics_{percentage}_{sampling_process}.npy")
        np.save(file=fname, arr=val_metrics, allow_pickle=True)

        # Update Variables
        # Min Training Loss
        if vanilla_avg_train_loss < min_train_loss:
            print(f"Train loss decreased from {min_train_loss} to {vanilla_avg_train_loss}.")
            min_train_loss = vanilla_avg_train_loss

        # DEBUG
        transf = transforms.ToPILImage()
        for j in range(len(W)):
            if W[j].sum() > 0:
                Xaug, X, labels, indices = train_loader.dataset[j]
                #print(len(X))
                #img1 = transf(X[j])
                #img1 = img1.save(f'/home/up201605633/Desktop/AL_debug/X-{epoch}-{j}.png')
                
                #img2 = transf(W[j])
                #plt.imshow(img2)
                #img2 = img2.save(f'/home/up201605633/Desktop/AL_debug/W-{epoch}-{j}.png')
                
        # Validation Loop
        print("Validation Phase")


        # Initialise lists to compute scores
        y_val_true = list()
        y_val_pred = list()


        # Running train loss
        run_val_loss = 0.0


        # Put model in evaluation mode
        model.eval()

        # Deactivate gradients
        with torch.no_grad():

            # Iterate through dataloader
            for batch_idx, (images, images_og, labels, indices) in enumerate(tqdm(val_loader)):

                # Move data data anda model to GPU (or not)
                images, labels = images.to(DEVICE), labels.to(DEVICE)
                model = model.to(DEVICE)

                # Forward pass: compute predicted outputs by passing inputs to the model
                logits = model(images)
                
                # Compute the batch loss
                # Using CrossEntropy w/ Softmax
                loss = LOSS(logits, labels)
                
                # Update batch losses
                run_val_loss += (loss.item() * images.size(0))

                # Concatenate lists
                y_val_true += list(labels.cpu().detach().numpy())
                
                # Using Softmax Activation
                # Apply Softmax on Logits and get the argmax to get the predicted labels
                s_logits = torch.nn.Softmax(dim=1)(logits)
                s_logits = torch.argmax(s_logits, dim=1)
                y_val_pred += list(s_logits.cpu().detach().numpy())

            # Compute Average Train Loss
            avg_val_loss = run_val_loss/len(val_loader.dataset)

            # Compute Training Accuracy
            val_acc = accuracy_score(y_true=y_val_true, y_pred=y_val_pred)
            val_recall = recall_score(y_true=y_val_true, y_pred=y_val_pred, average="weighted")
            val_precision = precision_score(y_true=y_val_true, y_pred=y_val_pred, average="weighted")
            val_f1 = f1_score(y_true=y_val_true, y_pred=y_val_pred, average="weighted")

            # Print Statistics
            print(f"Validation Loss: {avg_val_loss}\tValidation Accuracy: {val_acc}")
            # print(f"Validation Loss: {avg_val_loss}\tValidation Accuracy: {val_acc}\tValidation Recall: {val_recall}\tValidation Precision: {val_precision}\tValidation F1-Score: {val_f1}")

            # Append values to the arrays
            # Train Loss
            val_losses[epoch] = avg_val_loss
            # Save it to directory
            fname = os.path.join(history_dir, f"{model_name}_val_losses_{percentage}_{sampling_process}.npy")
            np.save(file=fname, arr=val_losses, allow_pickle=True)


            # Train Metrics
            # Acc
            val_metrics[epoch, 0] = val_acc
            # Recall
            val_metrics[epoch, 1] = val_recall
            # Precision
            val_metrics[epoch, 2] = val_precision
            # F1-Score
            val_metrics[epoch, 3] = val_f1
            # Save it to directory
            fname = os.path.join(history_dir, f"{model_name}_val_metrics_{percentage}_{sampling_process}.npy")
            np.save(file=fname, arr=val_metrics, allow_pickle=True)

            # Update Variables
            # Min validation loss and save if validation loss decreases
            if avg_val_loss < min_val_loss:
                print(f"Validation loss decreased from {min_val_loss} to {avg_val_loss}.")
                min_val_loss = avg_val_loss

                print("Saving best model on validation...")

                # Save checkpoint
                model_path = os.path.join(weights_dir, f"{model_name}_{percentage}p_{EPOCHS}e_{sampling_process}.pt")
                torch.save(model.state_dict(), model_path)

                print(f"Successfully saved at: {model_path}")
    

    model_path = os.path.join(weights_dir, f"{model_name}_{percentage}p_{EPOCHS}e_{sampling_process}.pt")
    torch.save(model.state_dict(), model_path)
    # Finish statement
    print("Finished.")
    return val_losses,train_losses,val_metrics,train_metrics
 
