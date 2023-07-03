import numpy as np

DEFAULT_EPSILON = 0.1

# Ordinary Least Squares
def least_squares(
        X_train, 
        y_train, 
        X_test, 
        *,
        fit_intercept=True,
        copy_X=True,
        n_jobs=None,
        positive=False):

    # import matplotlib.pyplot as plt
    # import numpy as np

    # Load the diabetes dataset
    # diabetes_X, diabetes_y = datasets.load_diabetes(return_X_y=True)

    # Use only one feature
    # diabetes_X = diabetes_X[:, np.newaxis, 2]

    # Split the data into training/testing sets
    # diabetes_X_train = diabetes_X[:-20]
    # diabetes_X_test = diabetes_X[-20:]

    # Split the targets into training/testing sets
    # diabetes_y_train = diabetes_y[:-20]
    # diabetes_y_test = diabetes_y[-20:]

    from sklearn.linear_model import LinearRegression
    # from sklearn.metrics import mean_squared_error, r2_score

    # Create linear regression object
    regr = LinearRegression(
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        n_jobs=n_jobs,
        positive=positive
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

    # The coefficients
    # print("Coefficients: \n", regr.coef_)
    # The mean squared error
    # print("Mean squared error: %.2f" % mean_squared_error(diabetes_y_test, diabetes_y_pred))
    # The coefficient of determination: 1 is perfect prediction
    # print("Coefficient of determination: %.2f" % r2_score(diabetes_y_test, diabetes_y_pred))

    # Plot outputs
    # plt.scatter(diabetes_X_test, diabetes_y_test, color="black")
    # plt.plot(diabetes_X_test, diabetes_y_pred, color="blue", linewidth=3)

    # plt.xticks(()) 
    # plt.yticks(())

    # plt.show()

# Ridge Regression and Classification
def ridge_regression(
        X_train, 
        y_train, 
        X_test,
        alpha=1.0,
        *,
        fit_intercept=True,
        copy_X=True,
        max_iter=None,
        tol=1e-4,
        solver="auto",
        positive=False,
        random_state=None):
    
    from sklearn.linear_model import Ridge

    # Create ridge regression object
    regr = Ridge(
        alpha=alpha,
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        max_iter=max_iter,
        tol=tol,
        solver=solver,
        positive=positive,
        random_state=random_state)

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def ridge_regression_cv(
        X_train, 
        y_train, 
        X_test,
        alphas=(0.1, 1.0, 10.0),
        *,
        fit_intercept=True,
        scoring=None,
        cv=None,
        gcv_mode=None,
        store_cv_values=False,
        alpha_per_target=False):
    
    from sklearn.linear_model import RidgeCV

    # Create ridge regression object
    regr = RidgeCV(
        alphas=alphas,
        fit_intercept=fit_intercept,
        scoring=scoring,
        cv=cv,
        gcv_mode=gcv_mode,
        store_cv_values=store_cv_values,
        alpha_per_target=alpha_per_target
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def ridge_classifier(
        X_train, 
        y_train, 
        X_test,
        alpha=1.0,
        *,
        fit_intercept=True,
        copy_X=True,
        max_iter=None,
        tol=1e-4,
        class_weight=None,
        solver="auto",
        positive=False,
        random_state=None):
    
    from sklearn.linear_model import RidgeClassifier

    # Create ridge classifier object
    regr = RidgeClassifier(
        alpha=alpha,
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        max_iter=max_iter,
        tol=tol,
        class_weight=class_weight,
        solver=solver,
        positive=positive,
        random_state=random_state
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr.coef_, y_pred

def ridge_classifier_cv(
        X_train, 
        y_train, 
        X_test,
        alphas=(0.1, 1.0, 10.0),
        *,
        fit_intercept=True,
        scoring=None,
        cv=None,
        class_weight=None,
        store_cv_values=False):
    
    from sklearn.linear_model import RidgeClassifierCV

    # Create ridge classifier with cross-validation object
    regr = RidgeClassifierCV(
        alphas=alphas,
        fit_intercept=fit_intercept,
        scoring=scoring,
        cv=cv,
        class_weight=class_weight,
        store_cv_values=store_cv_values
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Lasso
def lasso(
        X_train, 
        y_train, 
        X_test,
        alpha=1.0,
        *,
        fit_intercept=True,
        precompute=False,
        copy_X=True,
        max_iter=1000,
        tol=1e-4,
        warm_start=False,
        positive=False,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import Lasso

    # Create lasso regressor trained with L1 prior as regularizer object
    regr = Lasso(
        alpha=alpha,
        fit_intercept=fit_intercept,
        precompute=precompute,
        copy_X=copy_X,
        max_iter=max_iter,
        tol=tol,
        warm_start=warm_start,
        positive=positive,
        random_state=random_state,
        selection=selection
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def lasso_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        eps=1e-3,
        n_alphas=100,
        alphas=None,
        fit_intercept=True,
        precompute="auto",
        max_iter=1000,
        tol=1e-4,
        copy_X=True,
        cv=None,
        verbose=False,
        n_jobs=None,
        positive=False,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import LassoCV

    # Create lasso classifier with iterative fitting along a regularization path object
    regr = LassoCV(
        eps=eps,
        n_alphas=n_alphas,
        alphas=alphas,
        fit_intercept=fit_intercept,
        precompute=precompute,
        max_iter=max_iter,
        tol=tol,
        copy_X=copy_X,
        cv=cv,
        verbose=verbose,
        n_jobs=n_jobs,
        positive=positive,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# LARS Lasso
def lasso_lars(
        X_train, 
        y_train, 
        X_test,
        *,
        fit_intercept=True,
        verbose=False,
        normalize="deprecated",
        precompute="auto",
        max_iter=500,
        eps=np.finfo(float).eps,
        copy_X=True,
        fit_path=True,
        positive=False,
        jitter=None,
        random_state=None):
    
    from sklearn.linear_model import LassoLars

    # Create lasso classifier with Least Angle Regression object
    regr = LassoLars(
        fit_intercept=fit_intercept,
        verbose=verbose,
        normalize=normalize,
        precompute=precompute,
        max_iter=max_iter,
        eps=eps,
        copy_X=copy_X,
        fit_path=fit_path,
        positive=positive,
        jitter=jitter,
        random_state=random_state,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def lasso_lars_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        fit_intercept=True,
        verbose=False,
        max_iter=500,
        normalize="deprecated",
        precompute="auto",
        cv=None,
        max_n_alphas=1000,
        n_jobs=None,
        eps=np.finfo(float).eps,
        copy_X=True,
        positive=False):
    
    from sklearn.linear_model import LassoLarsCV

    # Create lasso regressor with cross-validation object
    regr = LassoLarsCV(
        fit_intercept=fit_intercept,
        verbose=verbose,
        max_iter=max_iter,
        normalize=normalize,
        precompute=precompute,
        cv=cv,
        max_n_alphas=max_n_alphas,
        n_jobs=n_jobs,
        eps=eps,
        copy_X=copy_X,
        positive=positive,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def lasso_lars_ic(
        X_train, 
        y_train, 
        X_test,
        *,
        fit_intercept=True,
        verbose=False,
        normalize="deprecated",
        precompute="auto",
        max_iter=500,
        eps=np.finfo(float).eps,
        copy_X=True,
        positive=False,
        noise_variance=None):
    
    from sklearn.linear_model import LassoLarsIC

    # Create lasso lars regressor using BIC or AIC object
    regr = LassoLarsIC(
        fit_intercept=fit_intercept,
        verbose=verbose,
        normalize=normalize,
        precompute=precompute,
        max_iter=max_iter,
        eps=eps,
        copy_X=copy_X,
        positive=positive,
        noise_variance=noise_variance,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Multi-task Lasso
def multi_task_lasso(
        X_train, 
        y_train, 
        X_test,
        alpha=1.0,
        *,
        fit_intercept=True,
        copy_X=True,
        max_iter=1000,
        tol=1e-4,
        warm_start=False,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import MultiTaskLasso

    # Create Multi-task Lasso model trained with L1/L2 mixed-norm as regularizer object
    regr = MultiTaskLasso(
        alpha=alpha,
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        max_iter=max_iter,
        tol=tol,
        warm_start=warm_start,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def multi_task_lasso_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        eps=1e-3,
        n_alphas=100,
        alphas=None,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-4,
        copy_X=True,
        cv=None,
        verbose=False,
        n_jobs=None,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import MultiTaskLassoCV

    # Create Multi-task Lasso model trained with L1/L2 mixed-norm as regularizer object
    regr = MultiTaskLassoCV(
        eps=eps,
        n_alphas=n_alphas,
        alphas=alphas,
        fit_intercept=fit_intercept,
        max_iter=max_iter,
        tol=tol,
        copy_X=copy_X,
        cv=cv,
        verbose=verbose,
        n_jobs=n_jobs,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Elastic-Net
def elastic_net(
        X_train, 
        y_train, 
        X_test,
        alpha=1.0,
        *,
        l1_ratio=0.5,
        fit_intercept=True,
        precompute=False,
        max_iter=1000,
        copy_X=True,
        tol=1e-4,
        warm_start=False,
        positive=False,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import ElasticNet

    # Create Linear regression with combined L1 and L2 priors as regularizer object
    regr = ElasticNet(
        alpha=alpha,
        l1_ratio=l1_ratio,
        fit_intercept=fit_intercept,
        precompute=precompute,
        max_iter=max_iter,
        copy_X=copy_X,
        tol=tol,
        warm_start=warm_start,
        positive=positive,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def elastic_net_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        l1_ratio=0.5,
        eps=1e-3,
        n_alphas=100,
        alphas=None,
        fit_intercept=True,
        precompute="auto",
        max_iter=1000,
        tol=1e-4,
        cv=None,
        copy_X=True,
        verbose=0,
        n_jobs=None,
        positive=False,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import ElasticNetCV

    # Create Elastic Net model with iterative fitting along a regularization path object
    regr = ElasticNetCV(
        l1_ratio=l1_ratio,
        eps=eps,
        n_alphas=n_alphas,
        alphas=alphas,
        fit_intercept=fit_intercept,
        precompute=precompute,
        max_iter=max_iter,
        tol=tol,
        cv=cv,
        copy_X=copy_X,
        verbose=verbose,
        n_jobs=n_jobs,
        positive=positive,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Multi-task Elastic-Net
def multi_task_elastic_net(
        X_train, 
        y_train, 
        X_test,
        alpha=1.0,
        *,
        l1_ratio=0.5,
        fit_intercept=True,
        copy_X=True,
        max_iter=1000,
        tol=1e-4,
        warm_start=False,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import MultiTaskElasticNet

    # Create Multi-task ElasticNet model trained with L1/L2 mixed-norm as regularizer object
    regr = MultiTaskElasticNet(
        alpha=alpha,
        l1_ratio=l1_ratio,
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        max_iter=max_iter,
        tol=tol,
        warm_start=warm_start,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def multi_task_elastic_net_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        l1_ratio=0.5,
        eps=1e-3,
        n_alphas=100,
        alphas=None,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-4,
        cv=None,
        copy_X=True,
        verbose=0,
        n_jobs=None,
        random_state=None,
        selection="cyclic"):
    
    from sklearn.linear_model import MultiTaskElasticNetCV

    # Create Multi-task L1/L2 ElasticNet with built-in cross-validation object
    regr = MultiTaskElasticNetCV(
        l1_ratio=l1_ratio,
        eps=eps,
        n_alphas=n_alphas,
        alphas=alphas,
        fit_intercept=fit_intercept,
        max_iter=max_iter,
        tol=tol,
        cv=cv,
        copy_X=copy_X,
        verbose=verbose,
        n_jobs=n_jobs,
        random_state=random_state,
        selection=selection,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Least Angle Regression
def lars(
        X_train, 
        y_train, 
        X_test,
        *,
        fit_intercept=True,
        verbose=False,
        normalize="deprecated",
        precompute="auto",
        n_nonzero_coefs=500,
        eps=np.finfo(float).eps,
        copy_X=True,
        fit_path=True,
        jitter=None,
        random_state=None):
    
    from sklearn.linear_model import Lars

    # Create Least Angle Regression model a.k.a. LAR object
    regr = Lars(
        fit_intercept=fit_intercept,
        verbose=verbose,
        normalize=normalize,
        precompute=precompute,
        n_nonzero_coefs=n_nonzero_coefs,
        eps=eps,
        copy_X=copy_X,
        fit_path=fit_path,
        jitter=jitter,
        random_state=random_state,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def lars_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        fit_intercept=True,
        verbose=False,
        max_iter=500,
        normalize="deprecated",
        precompute="auto",
        cv=None,
        max_n_alphas=1000,
        n_jobs=None,
        eps=np.finfo(float).eps,
        copy_X=True):
    
    from sklearn.linear_model import LarsCV

    # Create Cross-validated Least Angle Regression model object
    regr = LarsCV(
        fit_intercept=fit_intercept,
        verbose=verbose,
        max_iter=max_iter,
        normalize=normalize,
        precompute=precompute,
        cv=cv,
        max_n_alphas=max_n_alphas,
        n_jobs=n_jobs,
        eps=eps,
        copy_X=copy_X,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Orthogonal Matching Pursuit (OMP)
def orthogonal_matching_pursuit(
        X_train, 
        y_train, 
        X_test,
        *,
        n_nonzero_coefs=None,
        tol=None,
        fit_intercept=True,
        normalize="deprecated",
        precompute="auto"):
    
    from sklearn.linear_model import OrthogonalMatchingPursuit

    # Create Orthogonal Matching Pursuit model (OMP) object
    regr = OrthogonalMatchingPursuit(
        n_nonzero_coefs=n_nonzero_coefs,
        tol=tol,
        fit_intercept=fit_intercept,
        normalize=normalize,
        precompute=precompute,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def orthogonal_matching_pursuit_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        copy=True,
        fit_intercept=True,
        normalize="deprecated",
        max_iter=None,
        cv=None,
        n_jobs=None,
        verbose=False):
    
    from sklearn.linear_model import OrthogonalMatchingPursuitCV

    # Create Cross-validated Orthogonal Matching Pursuit model (OMP) object
    regr = OrthogonalMatchingPursuitCV(
        copy=copy,
        fit_intercept=fit_intercept,
        normalize=normalize,
        max_iter=max_iter,
        cv=cv,
        n_jobs=n_jobs,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Bayesian Regression
def bayesian_ridge_regression(
        X_train, 
        y_train, 
        X_test,
        *,
        n_iter=300,
        tol=1.0e-3,
        alpha_1=1.0e-6,
        alpha_2=1.0e-6,
        lambda_1=1.0e-6,
        lambda_2=1.0e-6,
        alpha_init=None,
        lambda_init=None,
        compute_score=False,
        fit_intercept=True,
        copy_X=True,
        verbose=False):
    
    from sklearn.linear_model import BayesianRidge

    # Create Bayesian ridge regression object
    regr = BayesianRidge(
        n_iter=n_iter,
        tol=tol,
        alpha_1=alpha_1,
        alpha_2=alpha_2,
        lambda_1=lambda_1,
        lambda_2=lambda_2,
        alpha_init=alpha_init,
        lambda_init=lambda_init,
        compute_score=compute_score,
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def bayesian_ard_regression(
        X_train, 
        y_train, 
        X_test,
        *,
        n_iter=300,
        tol=1.0e-3,
        alpha_1=1.0e-6,
        alpha_2=1.0e-6,
        lambda_1=1.0e-6,
        lambda_2=1.0e-6,
        compute_score=False,
        threshold_lambda=1.0e4,
        fit_intercept=True,
        copy_X=True,
        verbose=False):
    
    from sklearn.linear_model import ARDRegression

    # Create Bayesian ARD regression object
    regr = ARDRegression(
        n_iter=n_iter,
        tol=tol,
        alpha_1=alpha_1,
        alpha_2=alpha_2,
        lambda_1=lambda_1,
        lambda_2=lambda_2,
        compute_score=compute_score,
        threshold_lambda=threshold_lambda,
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Logistic Regression
def logistic_regression(
        X_train, 
        y_train, 
        X_test,
        penalty="l2",
        *,
        dual=False,
        tol=1e-4,
        C=1.0,
        fit_intercept=True,
        intercept_scaling=1,
        class_weight=None,
        random_state=None,
        solver="lbfgs",
        max_iter=100,
        multi_class="auto",
        verbose=0,
        warm_start=False,
        n_jobs=None,
        l1_ratio=None):
    
    from sklearn.linear_model import LogisticRegression

    # Create Logistic Regression (aka logit, MaxEnt) classifier object
    regr = LogisticRegression(
        penalty=penalty,
        dual=dual,
        tol=tol,
        C=C,
        fit_intercept=fit_intercept,
        intercept_scaling=intercept_scaling,
        class_weight=class_weight,
        random_state=random_state,
        solver=solver,
        max_iter=max_iter,
        multi_class=multi_class,
        verbose=verbose,
        warm_start=warm_start,
        n_jobs=n_jobs,
        l1_ratio=l1_ratio,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)
    y_prob_pred = regr.predict_proba(X_test)

    return regr, y_pred, y_prob_pred

def logistic_regression_cv(
        X_train, 
        y_train, 
        X_test,
        *,
        Cs=10,
        fit_intercept=True,
        cv=None,
        dual=False,
        penalty="l2",
        scoring=None,
        solver="lbfgs",
        tol=1e-4,
        max_iter=100,
        class_weight=None,
        n_jobs=None,
        verbose=0,
        refit=True,
        intercept_scaling=1.0,
        multi_class="auto",
        random_state=None,
        l1_ratios=None):
    
    from sklearn.linear_model import LogisticRegressionCV

    # Create Logistic Regression CV (aka logit, MaxEnt) classifier object
    regr = LogisticRegressionCV(
        Cs=Cs,
        fit_intercept=fit_intercept,
        cv=cv,
        dual=dual,
        penalty=penalty,
        scoring=scoring,
        solver=solver,
        tol=tol,
        max_iter=max_iter,
        class_weight=class_weight,
        n_jobs=n_jobs,
        verbose=verbose,
        refit=refit,
        intercept_scaling=intercept_scaling,
        multi_class=multi_class,
        random_state=random_state,
        l1_ratios=l1_ratios,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)
    y_prob_pred = regr.predict_proba(X_test)

    return regr, y_pred, y_prob_pred

# Generalized Linear Models
def tweedie_regressor(
        X_train, 
        y_train, 
        X_test,
        *,
        power=0.0,
        alpha=1.0,
        fit_intercept=True,
        link="auto",
        solver="lbfgs",
        max_iter=100,
        tol=1e-4,
        warm_start=False,
        verbose=0):
    
    from sklearn.linear_model import TweedieRegressor

    # Create Generalized Linear Model with a Tweedie distribution object
    regr = TweedieRegressor(
        power=power,
        alpha=alpha,
        fit_intercept=fit_intercept,
        link=link,
        solver=solver,
        max_iter=max_iter,
        tol=tol,
        warm_start=warm_start,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def poisson_regressor(
        X_train, 
        y_train, 
        X_test,
        *,
        alpha=1.0,
        fit_intercept=True,
        solver="lbfgs",
        max_iter=100,
        tol=1e-4,
        warm_start=False,
        verbose=0):
    
    from sklearn.linear_model import PoissonRegressor

    # Create Generalized Linear Model with a Poisson distribution object
    regr = PoissonRegressor(
        alpha=alpha,
        fit_intercept=fit_intercept,
        solver=solver,
        max_iter=max_iter,
        tol=tol,
        warm_start=warm_start,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def gamma_regressor(
        X_train, 
        y_train, 
        X_test,
        *,
        alpha=1.0,
        fit_intercept=True,
        solver="lbfgs",
        max_iter=100,
        tol=1e-4,
        warm_start=False,
        verbose=0):
    
    from sklearn.linear_model import GammaRegressor

    # Create Generalized Linear Model with a Poisson distribution object
    regr = GammaRegressor(
        alpha=alpha,
        fit_intercept=fit_intercept,
        solver=solver,
        max_iter=max_iter,
        tol=tol,
        warm_start=warm_start,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Stochastic Gradient Descent
def sgd_classifier(
        X_train, 
        y_train, 
        X_test,
        loss="hinge",
        *,
        penalty="l2",
        alpha=0.0001,
        l1_ratio=0.15,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-3,
        shuffle=True,
        verbose=0,
        epsilon=DEFAULT_EPSILON,
        n_jobs=None,
        random_state=None,
        learning_rate="optimal",
        eta0=0.0,
        power_t=0.5,
        early_stopping=False,
        validation_fraction=0.1,
        n_iter_no_change=5,
        class_weight=None,
        warm_start=False,
        average=False):
    
    from sklearn.preprocessing import StandardScaler
    from sklearn.linear_model import SGDClassifier
    from sklearn.pipeline import make_pipeline

    # Create Linear classifiers (SVM, logistic regression, etc.) with SGD training object
    clf = make_pipeline(StandardScaler(),
                        SGDClassifier(
                            loss=loss,
                            penalty=penalty,
                            alpha=alpha,
                            l1_ratio=l1_ratio,
                            fit_intercept=fit_intercept,
                            max_iter=max_iter,
                            tol=tol,
                            shuffle=shuffle,
                            verbose=verbose,
                            epsilon=epsilon,
                            n_jobs=n_jobs,
                            random_state=random_state,
                            learning_rate=learning_rate,
                            eta0=eta0,
                            power_t=power_t,
                            early_stopping=early_stopping,
                            validation_fraction=validation_fraction,
                            n_iter_no_change=n_iter_no_change,
                            class_weight=class_weight,
                            warm_start=warm_start,
                            average=average,
                        )
    )

    # Train the model using the training sets
    clf.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = clf.predict(X_test)

    return clf, y_pred

def sgd_regressor(
        X_train, 
        y_train, 
        X_test,
        loss="squared_error",
        *,
        penalty="l2",
        alpha=0.0001,
        l1_ratio=0.15,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-3,
        shuffle=True,
        verbose=0,
        epsilon=DEFAULT_EPSILON,
        random_state=None,
        learning_rate="invscaling",
        eta0=0.01,
        power_t=0.25,
        early_stopping=False,
        validation_fraction=0.1,
        n_iter_no_change=5,
        warm_start=False,
        average=False):
    
    from sklearn.preprocessing import StandardScaler
    from sklearn.linear_model import SGDRegressor
    from sklearn.pipeline import make_pipeline

    # Create Linear model fitted by minimizing a regularized empirical loss with SGD object
    regr = make_pipeline(StandardScaler(),
                        SGDRegressor(
                            loss=loss,
                            penalty=penalty,
                            alpha=alpha,
                            l1_ratio=l1_ratio,
                            fit_intercept=fit_intercept,
                            max_iter=max_iter,
                            tol=tol,
                            shuffle=shuffle,
                            verbose=verbose,
                            epsilon=epsilon,
                            random_state=random_state,
                            learning_rate=learning_rate,
                            eta0=eta0,
                            power_t=power_t,
                            early_stopping=early_stopping,
                            validation_fraction=validation_fraction,
                            n_iter_no_change=n_iter_no_change,
                            warm_start=warm_start,
                            average=average,
                        )
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Perceptron
def perceptron(
        X_train, 
        y_train, 
        X_test,
        *,
        penalty=None,
        alpha=0.0001,
        l1_ratio=0.15,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-3,
        shuffle=True,
        verbose=0,
        eta0=1.0,
        n_jobs=None,
        random_state=0,
        early_stopping=False,
        validation_fraction=0.1,
        n_iter_no_change=5,
        class_weight=None,
        warm_start=False):
    
    from sklearn.linear_model import Perceptron

    # Create Linear perceptron classifier object
    clf = Perceptron(
        penalty=penalty,
        alpha=alpha,
        l1_ratio=l1_ratio,
        fit_intercept=fit_intercept,
        max_iter=max_iter,
        tol=tol,
        shuffle=shuffle,
        verbose=verbose,
        eta0=eta0,
        n_jobs=n_jobs,
        random_state=random_state,
        early_stopping=early_stopping,
        validation_fraction=validation_fraction,
        n_iter_no_change=n_iter_no_change,
        class_weight=class_weight,
        warm_start=warm_start,
    )

    # Train the model using the training sets
    clf.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = clf.predict(X_test)

    return clf, y_pred

# Passive Aggressive Algorithms
def passive_aggressive_classifier(
        X_train, 
        y_train, 
        X_test,
        *,
        C=1.0,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-3,
        early_stopping=False,
        validation_fraction=0.1,
        n_iter_no_change=5,
        shuffle=True,
        verbose=0,
        loss="hinge",
        n_jobs=None,
        random_state=None,
        warm_start=False,
        class_weight=None,
        average=False):
    
    from sklearn.linear_model import PassiveAggressiveClassifier

    # Create Passive Aggressive Classifier object
    clf = PassiveAggressiveClassifier(
        C=C,
        fit_intercept=fit_intercept,
        max_iter=max_iter,
        tol=tol,
        early_stopping=early_stopping,
        validation_fraction=validation_fraction,
        n_iter_no_change=n_iter_no_change,
        shuffle=shuffle,
        verbose=verbose,
        loss=loss,
        n_jobs=n_jobs,
        random_state=random_state,
        warm_start=warm_start,
        class_weight=class_weight,
        average=average,
    )

    # Train the model using the training sets
    clf.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = clf.predict(X_test)

    return clf, y_pred

def passive_aggressive_regressor(
        X_train, 
        y_train, 
        X_test,
        *,
        C=1.0,
        fit_intercept=True,
        max_iter=1000,
        tol=1e-3,
        early_stopping=False,
        validation_fraction=0.1,
        n_iter_no_change=5,
        shuffle=True,
        verbose=0,
        loss="epsilon_insensitive",
        epsilon=DEFAULT_EPSILON,
        random_state=None,
        warm_start=False,
        average=False):
    
    from sklearn.linear_model import PassiveAggressiveRegressor

    # Create Passive Aggressive Regressor object
    regr = PassiveAggressiveRegressor(
        C=C,
        fit_intercept=fit_intercept,
        max_iter=max_iter,
        tol=tol,
        early_stopping=early_stopping,
        validation_fraction=validation_fraction,
        n_iter_no_change=n_iter_no_change,
        shuffle=shuffle,
        verbose=verbose,
        loss=loss,
        epsilon=epsilon,
        random_state=random_state,
        warm_start=warm_start,
        average=average,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Robustness regression
def huber_regression(
        X_train, 
        y_train, 
        X_test,
        *,
        epsilon=1.35,
        max_iter=100,
        alpha=0.0001,
        warm_start=False,
        fit_intercept=True,
        tol=1e-05):
    
    from sklearn.linear_model import HuberRegressor

    # Create Passive Aggressive Regressor object
    regr = HuberRegressor(
        epsilon=epsilon,
        max_iter=max_iter,
        alpha=alpha,
        warm_start=warm_start,
        fit_intercept=fit_intercept,
        tol=tol,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def ransac_regression(
        X_train, 
        y_train, 
        X_test,
        *,
        min_samples=None,
        residual_threshold=None,
        is_data_valid=None,
        is_model_valid=None,
        max_trials=100,
        max_skips=np.inf,
        stop_n_inliers=np.inf,
        stop_score=np.inf,
        stop_probability=0.99,
        loss="absolute_error",
        random_state=None,
        base_estimator="deprecated"):
    
    from sklearn.linear_model import RANSACRegressor

    # Create RANSAC (RANdom SAmple Consensus) algorithm object
    regr = RANSACRegressor(
        min_samples=min_samples,
        residual_threshold=residual_threshold,
        is_data_valid=is_data_valid,
        is_model_valid=is_model_valid,
        max_trials=max_trials,
        max_skips=max_skips,
        stop_n_inliers=stop_n_inliers,
        stop_score=stop_score,
        stop_probability=stop_probability,
        loss=loss,
        random_state=random_state,
        base_estimator=base_estimator,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

def theil_sen_regression(
        X_train, 
        y_train, 
        X_test,
        *,
        fit_intercept=True,
        copy_X=True,
        max_subpopulation=1e4,
        n_subsamples=None,
        max_iter=300,
        tol=1.0e-3,
        random_state=None,
        n_jobs=None,
        verbose=False,):
    
    from sklearn.linear_model import TheilSenRegressor

    # Create Theil-Sen Estimator: robust multivariate regression model object
    regr = TheilSenRegressor(
        fit_intercept=fit_intercept,
        copy_X=copy_X,
        max_subpopulation=max_subpopulation,
        n_subsamples=n_subsamples,
        max_iter=max_iter,
        tol=tol,
        random_state=random_state,
        n_jobs=n_jobs,
        verbose=verbose,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

# Quantile Regression
def quantile_regression(
        X_train, 
        y_train, 
        X_test,
        *,
        quantile=0.5,
        alpha=1.0,
        fit_intercept=True,
        solver="warn",
        solver_options=None):
    
    from sklearn.linear_model import QuantileRegressor

    # Create Linear regression model that predicts conditional quantile object
    regr = QuantileRegressor(
        quantile=quantile,
        alpha=alpha,
        fit_intercept=fit_intercept,
        solver=solver,
        solver_options=solver_options,
    )

    # Train the model using the training sets
    regr.fit(X_train, y_train)

    # Make predictions using the testing set
    y_pred = regr.predict(X_test)

    return regr, y_pred

