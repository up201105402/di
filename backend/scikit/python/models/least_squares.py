from model_functions import least_squares

regr, y_pred = least_squares(
    X_train, 
    y_train, 
    X_test,
    fit_intercept=<Fit_intercept>,
    copy_X=<Copy_X>,
    n_jobs=<N_jobs>,
    positive=<Positive>
)