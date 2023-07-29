import argparse

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--fit_intercept", action=argparse.BooleanOptionalAction)
    parser.add_argument("--string", type=str, required=False)
    args = parser.parse_args()

    print(args)

if __name__ == "__main__":
    main()