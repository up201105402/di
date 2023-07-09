from sklearn import datasets
import csv
import argparse
from pathlib import Path

def main():

    parser = argparse.ArgumentParser()
    parser.add_argument('-f1', "--orig_data_path", type=Path, required=True)
    parser.add_argument('-f2 ', "--orig_target_path", type=Path, required=True)
    parser.add_argument('-d', "--data_path", type=Path, required=True)
    parser.add_argument('-t', "--target_path", type=Path, required=True)
    parser.add_argument('-l1', "--lower_limit_data", type=int, required=False)
    parser.add_argument('-u1', "--upper_limit_data", type=int, required=False)
    parser.add_argument('-l2', "--lower_limit_target", type=int, required=False)
    parser.add_argument('-u2', "--upper_limit_target", type=int, required=False)
    args = parser.parse_args()

    with open(args.orig_data_path, newline='') as csvfile:
        X = list(csv.reader(csvfile))

    with open(args.orig_target_path, newline='') as csvfile:
        y = list(csv.reader(csvfile))

    with open(args.data_path, 'w') as csvfile:
        writer = csv.writer(csvfile)
        
        l1 = next(item for item in [args.lower_limit_data, 0] if item is not None)
        u1 = next(item for item in [args.upper_limit_data, len(X)] if item is not None)

        for row in X[l1:u1]:
            if not hasattr(row, '__iter__'):
                writer.writerow([ row ])
            else:
                writer.writerow(row)

    with open(args.target_path, 'w') as csvfile:
        writer = csv.writer(csvfile)

        l2 = next(item for item in [args.lower_limit_target, 0] if item is not None)
        u2 = next(item for item in [args.upper_limit_target, len(y)] if item is not None)

        for row in y[l2:u2]:
            if not hasattr(row, '__iter__'):
                writer.writerow([ row ])
            else:
                writer.writerow(row)

if __name__ == "__main__":
    main()