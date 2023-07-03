# datasets

def load_breast_cancer():
    from sklearn.datasets import load_breast_cancer

    # Load the Breast cancer wisconsin (diagnostic) dataset (classification)
    return load_breast_cancer(return_X_y=True)

def load_diabetes():
    from sklearn.datasets import load_diabetes

    # Load the diabetes dataset
    return load_diabetes(return_X_y=True)

def load_digits():
    from sklearn.datasets import load_digits

    # Load the digits dataset (classification)
    return load_digits(return_X_y=True)

def load_iris():
    from sklearn.datasets import load_iris

    # Load the iris dataset (classification)
    return load_iris(return_X_y=True)

def load_linerrud():
    from sklearn.datasets import load_linnerud

    # Load the Linnerrud dataset (classification)
    return load_linnerud(return_X_y=True)

def load_wine():
    from sklearn.datasets import load_wine

    # Load the Wine recognition dataset (classification)
    X, y = load_wine(return_X_y=True)