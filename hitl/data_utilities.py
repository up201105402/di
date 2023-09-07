# Imports
import os
import numpy as np
from tqdm import tqdm
from PIL import Image, ImageFile
import pandas as pd
from skimage.io import imread

# PyTorch Imports
import torch
from torch.utils.data import Dataset

import cv2
import csv
import shutil
import math

# Custom definitions
ImageFile.LOAD_TRUNCATED_IMAGES = True



# Function: Resize images based on new height
def resize_images(datapath, newpath, new_height=512):

    # Create new directories if needed
    if not os.path.exists(newpath):
        os.makedirs(newpath)
    

    # Go through the data
    for f in tqdm(os.listdir(datapath)):
        if(f.endswith(".jpg") or f.endswith('.png')):
            img = Image.open(os.path.join(datapath, f))
            w, h = img.size
            ratio = w / h
            new_w = int(np.ceil(new_height * ratio))
            new_img = img.resize((new_w, new_height), Image.ANTIALIAS)
            new_img.save(os.path.join(newpath, f))

    return



# NCI Dataset
# Class: NCI Dataset 
class NCI_Dataset(Dataset):
    def __init__(self, fold, path, transform=None, transform_orig=None, fraction=1):
        assert fold in ['train', 'test']
        # self.root = '/data/NCI/training/NHS'
        self.root = path
        self.files = [f for f in os.listdir(self.root) if f.endswith('_C1.jpg')]
        rand = np.random.RandomState(123)
        ix = rand.choice(len(self.files), len(self.files), False)
        ix = ix[:int(0.75*len(ix))] if fold == 'train' else ix[int(0.75*len(ix)):]
        self.files = [self.files[i] for i in ix]
        df = pd.read_excel(os.path.join(self.root, 'covariate_data_training_NHS.xls'), skiprows=2)
        self.classes = [df['WRST_HIST_AFTER'][df['IMAGE_ID'] == f].iloc[0] for f in self.files]
        self.transform = transform
        self.transform_orig = transform_orig
        
        # get desired fraction of data
        data_size = len(self.files)
        target_size = round(fraction * data_size)
        self.files = self.files[0:(target_size-1)]


    def __len__(self):
        return len(self.files)


    def __getitem__(self, idx):
        image = imread(os.path.join(self.root, self.files[idx]))
        label = self.classes[idx]
        label = 0 if label <= 0 else 1
        
        image = np.asarray(image)
        
        # Load image with PIL
        image = image_orig = Image.fromarray(image)

        # Apply transformation
        if self.transform:
            image = self.transform(image)
        if self.transform_orig:
            image_orig = self.transform_orig(image_orig)


        return image, image_orig, label, idx



# ISIC2017 Dataset
# Class: ISIC2017 Dataset
class ISIC17_Dataset(Dataset):
    def __init__(self, base_data_path, label_file, transform=None, transform_orig=None, fraction=1):
        """
        Args:
            base_data_path (string): Data directory.
            pickle_path (string): Path for pickle with annotations.
            transform (callable, optional): Optional transform to be applied
                on a sample.
        """
        
        # Init variables
        self.label_file = label_file
        self.base_data_path = base_data_path
        imgs_labels, self.labels_dict, self.nr_classes = self.map_images_and_labels(base_data_path, label_file)
        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]
        self.transform = transform
        self.transform_orig = transform_orig
        
        # get desired fraction of data
        data_size = len(imgs_labels)
        target_size = round(fraction * data_size)
        imgs_labels = imgs_labels[0:(target_size-1)]

        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]


        return 


    # Method: Map images and labels
    def map_images_and_labels(self, data_dir, label_file):
        
        # Get image_id and corresponding label from csv file
        labels = np.genfromtxt(os.path.join(data_dir, label_file), delimiter=',',encoding="utf8", dtype=None)
        labels = np.delete(labels,2,1)
        labels = np.delete(labels,0,0)
        
        # Images
        dir_files = os.listdir(data_dir)
        dir_imgs = [i for i in dir_files if i.split('.')[1]=='jpg']
        dir_imgs.sort()
        
        _labels_unique = np.unique(labels[:, 1])
        
        # Nr of Classes
        nr_classes = len(_labels_unique)

        # Create labels dictionary
        labels_dict = dict()
            
        for idx, _label in enumerate(_labels_unique):
            labels_dict[_label] = idx

        # Create img file name - image label array
        imgs_labels = np.column_stack((dir_imgs, labels[:,1]))
        
        return imgs_labels, labels_dict, nr_classes


    # Method: __len__
    def __len__(self):
        return len(self.images_paths)


    # Method: __getitem__
    def __getitem__(self, idx):
        if torch.is_tensor(idx):
            idx = idx.tolist()
        

        # Get images
        img_name = self.images_paths[idx]
        
        # Open image with PIL
        image = Image.open(os.path.join(self.base_data_path, img_name))
        #plt.imshow(image)
        
        # Perform transformations with Numpy Array
        image = np.asarray(image)
        
        #image = np.reshape(image, newshape=(image.shape[0], image.shape[1], 1))
        #image = np.concatenate((image, image, image), axis=2)

        # Load image with PIL
        image = image_orig = Image.fromarray(image)

        # Get labels
        label = self.labels_dict[self.images_labels[idx]]

        # Apply transformation
        if self.transform:
            image = self.transform(image)
        if self.transform_orig:
            image_orig = self.transform_orig(image_orig)

        #print(img_name, idx)
        return image, image_orig, label, idx


# APTOS2019 Dataset
# Class: APTOS2019 Dataset
class Aptos19_Dataset(Dataset):
    def __init__(self, base_data_path, label_file, transform=None, transform_orig=None, split= 'train', fraction=1):
        """
        Args:
            base_data_path (string): Data directory.
            pickle_path (string): Path for pickle with annotations.
            transform (callable, optional): Optional transform to be applied
                on a sample.
        """

        # Init variables
        self.label_file = label_file
        self.base_data_path = base_data_path
        imgs_labels, self.labels_dict, self.nr_classes = self.map_images_and_labels(base_data_path, label_file)

        # split train/test
        assert split in ['train', 'test']
        rand = np.random.RandomState(123)
        ix = rand.choice(len(imgs_labels), len(imgs_labels), False)
        if split == 'train':
            ix = ix[:int(len(ix)*0.8)]
        else:
            ix = ix[int(len(ix)*0.8):]
        imgs_labels = imgs_labels[ix]

        # get desired fraction of data
        data_size = len(imgs_labels)
        target_size = round(fraction * data_size)
        imgs_labels = imgs_labels[0:(target_size-1)]

        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]
        self.transform = transform
        self.transform_orig = transform_orig

        return


    # Method: Map images and labels
    def map_images_and_labels(self, data_dir, label_file):
        
        # Get image_id and corresponding label from csv file
        labels = np.genfromtxt(os.path.join(data_dir, label_file), delimiter=',',encoding="utf8", dtype=None)
        labels = np.delete(labels,0,0)

        # Images
        dir_files = os.listdir(data_dir)
        dir_imgs = [i for i in dir_files if i.split('.')[1]=='png']
        dir_imgs.sort()

        _labels_unique = np.unique(labels[:, 1])
        # Nr of Classes
        nr_classes = len(_labels_unique)

        # Create labels dictionary
        labels_dict = dict()
            
        for idx, _label in enumerate(_labels_unique):
            labels_dict[_label] = idx

        # Create img file name - image label array
        imgs_labels = np.column_stack((dir_imgs, labels[:,1]))

        return imgs_labels, labels_dict, nr_classes


    # Method: __len__
    def __len__(self):
        return len(self.images_paths)

    # Method: __getitem__
    def __getitem__(self, idx):
        if torch.is_tensor(idx):
            idx = idx.tolist()


        # Get images
        img_name = self.images_paths[idx]

        # Open image with PIL
        image = Image.open(os.path.join(self.base_data_path, img_name))
        #plt.imshow(image)

        # Perform transformations with Numpy Array
        image = np.asarray(image)

        #image = np.reshape(image, newshape=(image.shape[0], image.shape[1], 1))
        #image = np.concatenate((image, image, image), axis=2)

        # Load image with PIL
        image = image_orig = Image.fromarray(image)

        # Get labels
        label = self.labels_dict[self.images_labels[idx]]

        # Apply transformation
        if self.transform:
            image = self.transform(image)
        if self.transform_orig:
            image_orig = self.transform_orig(image_orig)

        #print(img_name, idx)
        return image, image_orig, label, idx


class ROSEYoutu_Dataset(Dataset):
    def __init__(self, base_data_path, split='train', transform=None, transform_orig=None, fraction=1):
        """
        Args:
            base_data_path (string): Data directory.
            pickle_path (string): Path for pickle with annotations.
            transform (callable, optional): Optional transform to be applied
                on a sample.
        """
        
        # Init variables
        assert split in ['train', 'test']
        self.label_file = os.path.join(base_data_path, split, 'labels.csv')
        self.base_data_path = os.path.join(base_data_path, split)

        # self.extract_frames_from_videos_folders(split, base_data_path)

        imgs_labels, self.labels_dict, self.nr_classes = self.map_images_and_labels(self.base_data_path, self.label_file)
        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]
        self.transform = transform
        self.transform_orig = transform_orig
        
        # get desired fraction of data
        data_size = len(imgs_labels)
        target_size = round(fraction * data_size)
        imgs_labels = imgs_labels[0:(target_size-1)]

        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]

        return 

    def extract_frames_from_videos_folders(self, split, base_data_path):

        folders = [int(folder) for folder in os.listdir(base_data_path) if os.path.isdir(os.path.join(base_data_path, folder)) and str.isdigit(folder)]
        folders.sort()

        folders = folders[:10] if split == 'train' else folders[-10:]

        target_dir = os.path.join(base_data_path, split)

        if not os.path.exists(target_dir):
            os.mkdir(target_dir)

        with open(f'{target_dir}/labels.csv', 'w', encoding='UTF-8') as f:
            csv_writer = csv.writer(f)
            csv_writer.writerow(['file', 'label'])

            for dir in folders:
                folder_path = os.path.join(base_data_path, str(dir))
                for file in os.listdir(folder_path):
                    video_frames = self.extract_frames_from_video(os.path.join(folder_path, file), 4)
                    filename, _ = os.path.splitext(file)
                    basename = os.path.basename(filename)

                    for i, frame in enumerate(video_frames):
                        frame_file_path = f"{target_dir}/{basename}_frame_{i}.png"
                        cv2.imwrite(frame_file_path, frame)
                        label = '0' if basename.split("_")[0] == 'G' else '1'
                        csv_writer.writerow([frame_file_path, label])

    def extract_frames_from_video(self, video_path, nr_frames):

        cap = cv2.VideoCapture(video_path)
        frame_count = cap.get(cv2.CAP_PROP_FRAME_COUNT)

        increment = frame_count // nr_frames
        offset = frame_count % nr_frames
        frames = []

        for i in range(nr_frames):
            cap.set(cv2.CAP_PROP_FRAME_COUNT, offset + (i*increment))
            ret, frame = cap.read()
            frames.append(frame)

        return frames

    # Method: Map images and labels
    def map_images_and_labels(self, data_dir, label_file):
        
        # Get image_id and corresponding label from csv file
        labels = np.genfromtxt(os.path.join(data_dir, label_file), delimiter=',',encoding="utf8", dtype=None)
        labels = np.delete(labels,0,0)
        
        # Images
        dir_files = os.listdir(data_dir)
        dir_imgs = [i for i in dir_files if i.split('.')[1]=='png']
        dir_imgs.sort()
        
        _labels_unique = np.unique(labels[:, 1])
        
        # Nr of Classes
        nr_classes = len(_labels_unique)

        # Create labels dictionary
        labels_dict = dict()
            
        for idx, _label in enumerate(_labels_unique):
            labels_dict[_label] = idx

        # Create img file name - image label array
        imgs_labels = np.column_stack((dir_imgs, labels[:,1]))
        
        return imgs_labels, labels_dict, nr_classes

    # Method: __len__
    def __len__(self):
        return len(self.images_paths)

    # Method: __getitem__
    def __getitem__(self, idx):
        if torch.is_tensor(idx):
            idx = idx.tolist()
        

        # Get images
        img_name = self.images_paths[idx]
        
        # Open image with PIL
        image = Image.open(os.path.join(self.base_data_path, img_name))
        #plt.imshow(image)
        
        # Perform transformations with Numpy Array
        image = np.asarray(image)
        
        #image = np.reshape(image, newshape=(image.shape[0], image.shape[1], 1))
        #image = np.concatenate((image, image, image), axis=2)

        # Load image with PIL
        image = image_orig = Image.fromarray(image)

        # Get labels
        label = self.labels_dict[self.images_labels[idx]]

        # Apply transformation
        if self.transform:
            image = self.transform(image)
        if self.transform_orig:
            image_orig = self.transform_orig(image_orig)

        #print(img_name, idx)
        return image, image_orig, label, idx

class PornographyXXX_Dataset(Dataset):
    def __init__(self, base_data_path, split='train', transform=None, transform_orig=None, fraction=1):
        """
        Args:
            base_data_path (string): Data directory.
            pickle_path (string): Path for pickle with annotations.
            transform (callable, optional): Optional transform to be applied
                on a sample.
        """
        
        # Init variables
        assert split in ['train', 'test']
        self.label_file = os.path.join(base_data_path, split, 'labels.csv')
        self.base_data_path = os.path.join(base_data_path, split)

        self.generate_label_file(split, base_data_path, fraction)

        imgs_labels, self.labels_dict, self.nr_classes = self.map_images_and_labels(self.base_data_path, self.label_file)
        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]
        self.transform = transform
        self.transform_orig = transform_orig

        self.images_paths, self.images_labels = imgs_labels[:, 0], imgs_labels[:, 1]

        return 

    def generate_label_file(self, split, base_data_path, fraction):

        folders = [folder for folder in os.listdir(base_data_path) if os.path.isdir(os.path.join(base_data_path, folder)) and folder != 'train' and folder != 'test']

        target_dir = os.path.join(base_data_path, split)

        if os.path.exists(target_dir):
            shutil.rmtree(target_dir)

        os.mkdir(target_dir)

        with open(f'{target_dir}/labels.csv', 'w', encoding='UTF-8') as f:
            csv_writer = csv.writer(f)
            csv_writer.writerow(['file', 'label'])

            for dir in folders:
                folder_path = os.path.join(base_data_path, dir)
                files = os.listdir(folder_path)
                file_count = len(files)
                fraction_size = math.floor(fraction * file_count)
                train_files = files[:fraction_size] if split == 'train' else files[-fraction_size:]

                for file in train_files:
                    shutil.copy(os.path.join(folder_path, file), target_dir)
                    label = '1' if dir == 'vPorn' else '0'
                    csv_writer.writerow([file, label])

    # Method: Map images and labels
    def map_images_and_labels(self, data_dir, label_file):
        
        # Get image_id and corresponding label from csv file
        labels = np.genfromtxt(label_file, delimiter=',', encoding="utf8", dtype=None, comments='/')
        labels = np.delete(labels,0,0)
        
        # Images
        dir_files = os.listdir(data_dir)
        dir_imgs = [i for i in dir_files if i.split('.')[1]=='jpg']
        dir_imgs.sort()
        
        _labels_unique = np.unique(labels[:, 1])
        
        # Nr of Classes
        nr_classes = len(_labels_unique)

        # Create labels dictionary
        labels_dict = dict()
            
        for idx, _label in enumerate(_labels_unique):
            labels_dict[_label] = idx

        # Create img file name - image label array
        imgs_labels = np.column_stack((dir_imgs, labels[:,1]))
        
        return imgs_labels, labels_dict, nr_classes


    # Method: __len__
    def __len__(self):
        return len(self.images_paths)


    # Method: __getitem__
    def __getitem__(self, idx):
        if torch.is_tensor(idx):
            idx = idx.tolist()
        

        # Get images
        img_name = self.images_paths[idx]
        
        # Open image with PIL
        image = Image.open(os.path.join(self.base_data_path, img_name))
        #plt.imshow(image)
        
        # Perform transformations with Numpy Array
        image = np.asarray(image)
        
        #image = np.reshape(image, newshape=(image.shape[0], image.shape[1], 1))
        #image = np.concatenate((image, image, image), axis=2)

        # Load image with PIL
        image = image_orig = Image.fromarray(image)

        # Get labels
        label = self.labels_dict[self.images_labels[idx]]

        # Apply transformation
        if self.transform:
            image = self.transform(image)
        if self.transform_orig:
            image_orig = self.transform_orig(image_orig)

        #print(img_name, idx)
        return image, image_orig, label, idx
