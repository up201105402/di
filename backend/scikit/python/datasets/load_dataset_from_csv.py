import pandas as pd
import csv
# from sklearn.linear_model import LinearRegression

# Read the dataset from the file
dataset = pd.read_csv(file_path)

# Split the dataset into features and target
X = dataset.iloc[:, :-1].values
y = dataset.iloc[:, -1].values

# Train the linear regression model
# regressor = LinearRegression()
# regressor.fit(X, y)

# Return the trained model
# return regressor