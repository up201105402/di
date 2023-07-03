from sklearn import datasets

# Load the Wine recognition dataset (classification)
X, y = datasets.load_wine(return_X_y=True)

print("Wine recognition dataset loaded!")