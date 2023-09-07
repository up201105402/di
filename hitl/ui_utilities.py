# Imports
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.patches import Rectangle

# PyTorch Imports
import torch
import torch.nn as nn


def matchSelectedRects(selectedRects, rects):
    indexes = []

    for i, rect in enumerate(rects):
        for selectedRect in selectedRects:
           if [j for j in selectedRect] == [t.item() for t in rect]:
               indexes.append(i)
    
    return {rects[i] for i in indexes}


# Function: See if a point is inside a rectangle
def point_inside_rect(pt, rect):
    return rect[0] <= pt[0] <= rect[2] and rect[1] <= pt[1] <= rect[3]



# Class: Choose rectangles (UI based in Matplotlib)
class ChooseRectangles:
    def __init__(self, img, rects, image_resize_factor=1, edgecolor='red'):
        self.img = img
        self.image_resize_factor = image_resize_factor
        self.rects = rects
        self.selected = set()
        self.edgecolor = edgecolor

    def draw_trimmed(self):
        plt.clf()
        fig, ax = plt.subplots(figsize=tuple([x/(100/self.image_resize_factor) for x in self.img.shape[:2]]), dpi=100)
        fig.subplots_adjust(0,0,1,1)
        ax.imshow(self.img)
        for i, (x1, y1, x2, y2) in enumerate(self.rects):
            if i in self.selected:
                plt.gca().add_patch(Rectangle(
                    (x1, y1), x2-x1, y2-y1,
                    facecolor=self.edgecolor, edgecolor='none', alpha=0.4))
            plt.gca().add_patch(Rectangle(
                (x1, y1), x2-x1, y2-y1,
                facecolor='none', edgecolor=self.edgecolor, lw=3))
        plt.draw()
    
    def draw(self):
        plt.clf()
        plt.imshow(self.img)
        for i, (x1, y1, x2, y2) in enumerate(self.rects):
            if i in self.selected:
                plt.gca().add_patch(Rectangle(
                    (x1, y1), x2-x1, y2-y1,
                    facecolor=self.edgecolor, edgecolor='none', alpha=0.4))
            plt.gca().add_patch(Rectangle(
                (x1, y1), x2-x1, y2-y1,
                facecolor='none', edgecolor=self.edgecolor, lw=3))
        # plt.connect('button_press_event', self.button_press)
        # plt.connect('key_press_event', self.key_press)
        plt.draw()
        plt.savefig('figure.png')


    # Method: Keyboard interface
    def key_press(self, event):
        if event.key in ['enter', 'escape']:
            plt.close()


    # Method: Button interface
    def button_press(self, event):
        if event.button == 1 and event.inaxes:
            pt = event.xdata, event.ydata
            i = [i for i, rect in enumerate(self.rects) if point_inside_rect(pt, rect)]
            if len(i) > 0:
                i = i[0]
                self.selected.discard(i) if i in self.selected else self.selected.add(i)
                self.draw()


    # Method: Get selected rectangles
    def get_selected_rectangles(self):
        return {self.rects[i] for i in self.selected}



# Class: Generate rectangles (UI based on Matplotlib)
class GenerateRectangles:
    def __init__(self, img, size, stride, nr_rects):
        self.img = img
        self.size = size
        self.stride = stride
        self.nr_rects = nr_rects


    # Method: Get ranked patches
    def get_ranked_patches(self):
        avg = nn.AvgPool2d(self.size, self.size)
        avg_patches = avg(self.img[None])

        print('avg_patches:', avg_patches.shape)

        # Sort patches
        ranks = torch.argsort(avg_patches.flatten()).reshape((avg_patches.shape))
        print('ranks shape:', ranks.shape)
        
        _, yi, xi = torch.where(ranks < self.nr_rects)


        rects = [(x*self.size, y*self.size, (x+1)*self.size, (y+1)*self.size) for y, x in zip(yi, xi)]

        return rects



# Function: Image show function
def imshow(img ,transpose = True):
    img = img / 2 + 0.5     # unnormalize
    npimg = img.numpy()
    plt.imshow(np.transpose(img, (1, 2, 0)))
    plt.show()



# Function: Get Oracle Feedback (UI)
def GetOracleFeedback(image, label, idx, model_attributions, pred, rectSize, rectStride, nr_rects, epoch_number, query_nr, image_resize_factor, epoch_dir):
    rectGenerator = GenerateRectangles(model_attributions, size=rectSize, stride=rectStride, nr_rects=nr_rects)
    rects = rectGenerator.get_ranked_patches()
    image = image / 2 + 0.5     # unnormalize
    #npimg = image.numpy()
    image = np.transpose(image, (1, 2, 0))
    
    ui = ChooseRectangles(image, rects, image_resize_factor)
    # ui.draw()
    # plt.title(f"True Label: {label}  Prediction: {pred}  Image idx: {idx}")
    # plt.show()
    # print(ui.selected)
    selected_rects = ui.get_selected_rectangles()

    ui.draw_trimmed()
    plt.show()
    plt.savefig(f"{epoch_dir}/query_{query_nr}_image.png")
    plt.close('all')
    
    torch.save(rects, f"{epoch_dir}/rects_{query_nr}.pt")

    pyRects = []
    for i, rect in enumerate(rects):
        indexedRect = [i.item() for i in rect]
        pyRects.append(indexedRect)
    
    np.savetxt(f"{epoch_dir}/query_{query_nr}_rects.csv", np.array(pyRects), fmt="%d", delimiter=",")

    return ui.selected, selected_rects
