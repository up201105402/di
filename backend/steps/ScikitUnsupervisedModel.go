package steps

import (
	"bytes"
	"di/model"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var ScikitUnsupervisedModelTypes = []string{
	"leastSquares",
	"ridgeRegression",
	"ridgeRegressionCV",
	"ridgeClassifier",
	"ridgeClassifierCV",
	"lasso",
	"lassoCV",
	"lassoLars",
	"lassoLarsCV",
	"lassoLarsIC",
	"multiTaskLasso",
	"multiTaskLassoCV",
	"elasticNet",
	"elasticNetCV",
	"multiTaskElasticNet",
	"multiTaskElasticNetCV",
	"lars",
	"larsCV",
	"omp",
	"ompCV",
	"bayesianRidge",
	"bayesianARD",
	"logisticRegression",
	"logisticRegressionCV",
	"tweedieRegressor",
	"poissonRegressor",
	"gammaRegressor",
	"sgdClassifier",
	"sgdRegressor",
	"perceptron",
	"passiveAgressiveClassifier",
	"passiveAgressiveRegressor",
	"huberRegression",
	"ransacRegression",
	"theilSenRegression",
	"quantileRegression",
}

type ScikitUnsupervisedModel struct {
	ID          int
	PipelineID  uint
	RunID       uint
	Model       string
	Name        string
	IsFirstStep bool
	DataConfig  model.StepDataConfig
}

func (step ScikitUnsupervisedModel) GetID() int {
	return int(step.ID)
}

func (step ScikitUnsupervisedModel) GetName() string {
	return step.Name
}

func (step ScikitUnsupervisedModel) GetIsFirstStep() bool {
	return step.IsFirstStep
}

func (step *ScikitUnsupervisedModel) SetData(stepDescription model.NodeDescription) error {
	step.ID, _ = strconv.Atoi(stepDescription.ID)
	step.Name = stepDescription.Data.NameAndType.Name
	step.IsFirstStep = stepDescription.Data.NameAndType.IsFirstStep
	step.Model = stepDescription.Data.Model.String
	step.DataConfig = stepDescription.Data.StepConfig
	return nil
}

func (step *ScikitUnsupervisedModel) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID
	return nil
}

func (step *ScikitUnsupervisedModel) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitUnsupervisedModel) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitUnsupervisedModel) GetRunID() uint {
	return step.RunID
}

func (step ScikitUnsupervisedModel) Execute(logFile *os.File, I18n *i18n.Localizer) error {

	runLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	var args []string

	args, err := step.appendArgs(args, I18n, runLogger)

	if err != nil {
		return err
	}

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})

		runLogger.Println(errMessage)
		return errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"

	cmd := exec.Command("python3", args...)
	cmd.Dir = currentPipelineWorkDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmdErr := cmd.Run()
	runLogger.Println(stderr.String())
	runLogger.Println(stdout.String())

	return cmdErr
}

func (step ScikitUnsupervisedModel) appendArgs(args []string, I18n *i18n.Localizer, runLogger *log.Logger) ([]string, error) {

	scikitSnippetsDir, exists := os.LookupEnv("SCIKIT_SNIPPETS_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "SCIKIT_SNIPPETS_DIR",
			},
			PluralCount: 1,
		})

		runLogger.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	scikitSnippetsDir = scikitSnippetsDir + "models/"

	args = append(args, scikitSnippetsDir+"linear_models.py")
	args = append(args, "--model")
	args = append(args, step.Model)

	pipelinesWorkDir, exists := os.LookupEnv("PIPELINES_WORK_DIR")

	if !exists {
		errMessage := I18n.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "env.variable.find.failed",
			TemplateData: map[string]interface{}{
				"Name": "PIPELINES_WORK_DIR",
			},
			PluralCount: 1,
		})

		runLogger.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	currentPipelineWorkDir := pipelinesWorkDir + "/" + fmt.Sprint(step.PipelineID) + "/" + fmt.Sprint(step.RunID) + "/"

	args = append(args, "--train_data_path")
	args = append(args, currentPipelineWorkDir+"filtered_training_data.csv")

	args = append(args, "--train_target_path")
	args = append(args, currentPipelineWorkDir+"filtered_training_target.csv")

	args = append(args, "--testing_data_path")
	args = append(args, currentPipelineWorkDir+"filtered_testing_data.csv")

	if step.DataConfig.Fit_intercept.Valid {
		args = append(args, "--fit_intercept")
	}

	if step.DataConfig.Copy_X.Valid && step.DataConfig.Copy_X.Bool {
		args = append(args, "--copy_X")
	}

	if step.DataConfig.N_jobs.Valid {
		args = append(args, "--n_jobs")
		args = append(args, string(step.DataConfig.N_jobs.Int64))
	}

	if step.DataConfig.Positive.Valid && step.DataConfig.Positive.Bool {
		args = append(args, "--positive")
	}

	if step.DataConfig.Alpha.Valid {
		args = append(args, "--alpha")
		args = append(args, step.DataConfig.Alpha.String)
	}

	if step.DataConfig.Max_iter.Valid {
		args = append(args, "--max_iter")
		args = append(args, step.DataConfig.Max_iter.String)
	}

	if step.DataConfig.Tol.Valid {
		args = append(args, "--tol")
		args = append(args, step.DataConfig.Tol.String)
	}

	if step.DataConfig.Solver.Valid {
		args = append(args, "--solver")
		args = append(args, step.DataConfig.Solver.String)
	}

	if step.DataConfig.Random_state.Valid {
		args = append(args, "--random_state")
		args = append(args, step.DataConfig.Random_state.String)
	}

	if step.DataConfig.Alphas.Valid {
		args = append(args, "--alphas")
		args = append(args, step.DataConfig.Alphas.String)
	}

	if step.DataConfig.Class_weight.Valid {
		args = append(args, "--class_weight")
		args = append(args, step.DataConfig.Class_weight.String)
	}

	if step.DataConfig.Scoring.Valid {
		args = append(args, "--scoring")
		args = append(args, step.DataConfig.Scoring.String)
	}

	if step.DataConfig.Cv.Valid {
		args = append(args, "--cv")
		args = append(args, step.DataConfig.Cv.String)
	}

	if step.DataConfig.Store_cv_values.Valid && step.DataConfig.Store_cv_values.Bool {
		args = append(args, "--store_cv_values")
	}

	if step.DataConfig.Precompute.Valid && step.DataConfig.Precompute.Bool {
		args = append(args, "--precompute")
	}

	if step.DataConfig.Warm_start.Valid {
		args = append(args, "--warm_start")
		args = append(args, step.DataConfig.Warm_start.String)
	}

	if step.DataConfig.Selection.Valid {
		args = append(args, "--selection")
		args = append(args, step.DataConfig.Selection.String)
	}

	if step.DataConfig.Eps.Valid {
		args = append(args, "--eps")
		args = append(args, step.DataConfig.Eps.String)
	}

	if step.DataConfig.N_alphas.Valid {
		args = append(args, "--n_alphas")
		args = append(args, step.DataConfig.N_alphas.String)
	}

	if step.DataConfig.Verbose.Valid && step.DataConfig.Verbose.Bool {
		args = append(args, "--verbose")
	}

	if step.DataConfig.Fit_path.Valid && step.DataConfig.Fit_path.Bool {
		args = append(args, "--fit_path")
	}

	if step.DataConfig.Jitter.Valid {
		args = append(args, "--jitter")
		args = append(args, step.DataConfig.Jitter.String)
	}

	if step.DataConfig.Max_n_alphas.Valid {
		args = append(args, "--max_n_alphas")
		args = append(args, step.DataConfig.Max_n_alphas.String)
	}

	if step.DataConfig.Criterion.Valid {
		args = append(args, "--criterion")
		args = append(args, step.DataConfig.Criterion.String)
	}

	if step.DataConfig.Noise_variance.Valid {
		args = append(args, "--noise_variance")
		args = append(args, step.DataConfig.Noise_variance.String)
	}

	if step.DataConfig.L1_ratio.Valid {
		args = append(args, "--l1_ratio")
		args = append(args, step.DataConfig.L1_ratio.String)
	}

	if step.DataConfig.N_nonzero_coefs.Valid {
		args = append(args, "--n_nonzero_coefs")
		args = append(args, step.DataConfig.N_nonzero_coefs.String)
	}

	if step.DataConfig.Copy.Valid && step.DataConfig.Copy.Bool {
		args = append(args, "--copy")
	}

	if step.DataConfig.N_iter.Valid {
		args = append(args, "--n_iter")
		args = append(args, step.DataConfig.N_iter.String)
	}

	if step.DataConfig.Alpha_1.Valid {
		args = append(args, "--alpha_1")
		args = append(args, step.DataConfig.Alpha_1.String)
	}

	if step.DataConfig.Alpha_2.Valid {
		args = append(args, "--alpha_2")
		args = append(args, step.DataConfig.Alpha_2.String)
	}

	if step.DataConfig.Lambda_1.Valid {
		args = append(args, "--lambda_1")
		args = append(args, step.DataConfig.Lambda_1.String)
	}

	if step.DataConfig.Lambda_2.Valid {
		args = append(args, "--lambda_2")
		args = append(args, step.DataConfig.Lambda_2.String)
	}

	if step.DataConfig.Alpha_init.Valid {
		args = append(args, "--alpha_init")
		args = append(args, step.DataConfig.Alpha_init.String)
	}

	if step.DataConfig.Lambda_init.Valid {
		args = append(args, "--lambda_init")
		args = append(args, step.DataConfig.Lambda_init.String)
	}

	if step.DataConfig.Compute_score.Valid && step.DataConfig.Compute_score.Bool {
		args = append(args, "--compute_score")
	}

	if step.DataConfig.Threshold_lambda.Valid {
		args = append(args, "--threshold_lambda")
		args = append(args, step.DataConfig.Threshold_lambda.String)
	}

	if step.DataConfig.Penalty.Valid {
		args = append(args, "--penalty")
		args = append(args, step.DataConfig.Penalty.String)
	}

	if step.DataConfig.Dual.Valid && step.DataConfig.Dual.Bool {
		args = append(args, "--dual")
	}

	if step.DataConfig.C.Valid {
		args = append(args, "--C")
		args = append(args, step.DataConfig.C.String)
	}

	if step.DataConfig.Intercept_scaling.Valid {
		args = append(args, "--intercept_scaling")
		args = append(args, step.DataConfig.Intercept_scaling.String)
	}

	if step.DataConfig.Multi_class.Valid {
		args = append(args, "--multi_class")
		args = append(args, step.DataConfig.Multi_class.String)
	}

	if step.DataConfig.Cs.Valid {
		args = append(args, "--Cs")
		args = append(args, step.DataConfig.Cs.String)
	}

	if step.DataConfig.Refit.Valid && step.DataConfig.Refit.Bool {
		args = append(args, "--refit")
	}

	if step.DataConfig.L1_ratios.Valid {
		args = append(args, "--l1_ratios")
		args = append(args, step.DataConfig.L1_ratios.String)
	}

	if step.DataConfig.Power.Valid {
		args = append(args, "--power")
		args = append(args, step.DataConfig.Power.String)
	}

	if step.DataConfig.Link.Valid {
		args = append(args, "--link")
		args = append(args, step.DataConfig.Link.String)
	}

	if step.DataConfig.Loss.Valid {
		args = append(args, "--loss")
		args = append(args, step.DataConfig.Loss.String)
	}

	if step.DataConfig.Shuffle.Valid && step.DataConfig.Shuffle.Bool {
		args = append(args, "--shuffle")
	}

	if step.DataConfig.Epsilon.Valid {
		args = append(args, "--epsilon")
		args = append(args, step.DataConfig.Epsilon.String)
	}

	if step.DataConfig.Learning_rate.Valid {
		args = append(args, "--learning_rate")
		args = append(args, step.DataConfig.Learning_rate.String)
	}

	if step.DataConfig.Eta0.Valid {
		args = append(args, "--eta0")
		args = append(args, step.DataConfig.Eta0.String)
	}

	if step.DataConfig.Power_t.Valid {
		args = append(args, "--power_t")
		args = append(args, step.DataConfig.Power_t.String)
	}

	if step.DataConfig.Early_stopping.Valid && step.DataConfig.Early_stopping.Bool {
		args = append(args, "--early_stopping")
	}

	if step.DataConfig.Validation_fraction.Valid {
		args = append(args, "--validation_fraction")
		args = append(args, step.DataConfig.Validation_fraction.String)
	}

	if step.DataConfig.N_iter_no_change.Valid {
		args = append(args, "--n_iter_no_change")
		args = append(args, step.DataConfig.N_iter_no_change.String)
	}

	if step.DataConfig.Average.Valid && step.DataConfig.Average.Bool {
		args = append(args, "--average")
	}

	if step.DataConfig.Max_subpopulation.Valid {
		args = append(args, "--max_subpopulation")
		args = append(args, step.DataConfig.Max_subpopulation.String)
	}

	if step.DataConfig.N_subsamples.Valid {
		args = append(args, "--n_subsamples")
		args = append(args, step.DataConfig.N_subsamples.String)
	}

	if step.DataConfig.Quantile.Valid {
		args = append(args, "--quantile")
		args = append(args, step.DataConfig.Quantile.String)
	}

	if step.DataConfig.Solver_options.Valid {
		args = append(args, "--solver_options")
		args = append(args, step.DataConfig.Solver_options.String)
	}

	return args, nil
}
