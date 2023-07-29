from sklearn import datasets
import csv
import argparse
import numpy as np
from pathlib import Path

def main():

    parser = argparse.ArgumentParser()
    parser.add_argument('-d', "--data_path", type=Path, required=True)
    parser.add_argument('-t', "--target_path", type=Path, required=False)
    parser.add_argument('-l1', "--lower_limit_data", type=Path, required=False)
    parser.add_argument('-u1', "--upper_limit_data", type=Path, required=False)
    parser.add_argument('-l2', "--lower_limit_target", type=Path, required=False)
    parser.add_argument('-u2', "--upper_limit_target", type=Path, required=False)
    args = parser.parse_args()

    X, y = datasets.load_diabetes(return_X_y=True)

    with open(args.data_path, 'wb') as file:
        
        l1 = next(item for item in [args.lower_limit_data, 0] if item is not None)
        u1 = next(item for item in [args.upper_limit_data, len(X)] if item is not None)

        np.save(file, X[l1:u1-1])

    if args.target_path:
        with open(args.target_path, 'wb') as file:

            l2 = next(item for item in [args.lower_limit_target, 0] if item is not None)
            u2 = next(item for item in [args.upper_limit_target, len(y)] if item is not None)

            np.save(file, y[l2:u2-1])

if __name__ == "__main__":
    main()