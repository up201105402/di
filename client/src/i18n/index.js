import { createI18n } from 'vue-i18n';

export const messages = {
    en: {
        pages: {
            home: {
                name: 'Home',
            },
            dashboard: {
                completedTasks: "Completed Tasks",
                dailySales: "Daily Sales",
                performance: "Performance",
                simpleTable: "Simple Table",
                totalShipments: "Total Shipments",
                chartCategories: [
                    "Accounts",
                    "Purchases",
                    "Sessions"
                ],
                tasks: "Tasks({count})",
                today: "Today",
                dropdown: {
                    action: "Action",
                    anotherAction: "Another Action",
                    somethingElse: "Something else"
                },
                usersTable: {
                    title: "Simple Table",
                    columns: [
                        "Name",
                        "Country",
                        "City",
                        "Salary"
                    ],
                    data: [
                        {
                            id: 1,
                            name: "Dakota Rice",
                            salary: "$36.738",
                            country: "Niger",
                            city: "Oud-Turnhout"
                        },
                        {
                            id: 2,
                            name: "Minerva Hooper",
                            salary: "$23,789",
                            country: "Curaçao",
                            city: "Sinaai-Waas"
                        },
                        {
                            id: 3,
                            name: "Sage Rodriguez",
                            salary: "$56,142",
                            country: "Netherlands",
                            city: "Baileux"
                        },
                        {
                            id: 4,
                            name: "Philip Chaney",
                            salary: "$38,735",
                            country: "Korea, South",
                            city: "Overland Park"
                        },
                        {
                            id: 5,
                            name: "Doris Greene",
                            salary: "$63,542",
                            country: "Malawi",
                            city: "Feldkirchen in Kärnten"
                        },
                        {
                            id: 6,
                            name: "Mason Porter",
                            salary: "$98,615",
                            country: "Chile",
                            city: "Gloucester"
                        },
                        {
                            id: 7,
                            name: "Jon Porter",
                            salary: "$78,615",
                            country: "Portugal",
                            city: "Gloucester"
                        }
                    ]
                },
                taskList: [
                    {
                        title: "Update the Documentation",
                        description: "Dwuamish Head, Seattle, WA 8:47 AM",
                        done: false
                    },
                    {
                        title: "GDPR Compliance",
                        description: "The GDPR is a regulation that requires businesses to protect the personal data and privacy of Europe citizens for transactions that occur within EU member states.",
                        done: true
                    },
                    {
                        title: "Solve the issues",
                        description: "Fifty percent of all respondents said they would be more likely to shop at a company",
                        done: false
                    },
                    {
                        title: "Release v2.0.0",
                        description: "Ra Ave SW, Seattle, WA 98116, SUA 11:19 AM",
                        done: false
                    },
                    {
                        title: "Export the processed files",
                        description: "The report also shows that consumers will not easily forgive a company once a breach exposing their personal data occurs.",
                        done: false
                    },
                    {
                        title: "Arival at export process",
                        description: "Capitol Hill, Seattle, WA 12:34 AM",
                        done: false
                    }
                ]
            },
            pipelines: {
                name: "Pipelines",
                header: "Pipelines",
                table: {
                    headers: {
                        name: 'Name', 
                        modified: 'Modified',
                        created: 'Created',
                    },
                },
                dialog: {
                    create: {
                        header: 'Create Pipeline',
                        name: {
                            label: 'Name',
                            help: 'Please enter the pipeline name',
                        }
                    },
                    delete: {
                        header: 'Confirm Delete',
                        body: 'This will permanently delete this pipeline!',
                    },
                },
                edit: {
                    header: "Pipeline {id}",
                    dialog: {
                        create: {
                            header: 'Create {name} Step',
                        },
                        edit: {
                            header: 'Edit {name} Step',
                        },
                        nameAndType: {
                            label: 'Name and Type',
                            name: 'Step Name',
                            isFirstStep: 'Is First Step?',
                            scriptType: {
                                label: 'Script Type',
                                options: {
                                    inline: 'Inline',
                                    file: 'File',
                                }
                            },
                            dataset: {
                                label: 'Dataset',
                                options: {
                                    scikitBreastCancer: 'Breast Cancer Dataset', 
                                    scikitDiabetes: 'Diabetes Dataset', 
                                    scikitDigits: 'Digits Dataset', 
                                    scikitIris: 'Iris Dataset', 
                                    scikitLinerrud: 'Linnerud Dataset', 
                                    scikitWine: 'Wine Dataset', 
                                    scikitLoadFile: 'Load Daset From File', 
                                }
                            },
                        },
                        stepConfig: {
                            label: 'Step Config',
                            repoUrl: 'URL',
                            scriptEditor: 'Script Editor',
                            scriptFile: {
                                label: 'File',
                                button: 'File Upload',
                            },
                            dataFilePath: 'Data File Path',
                            targetFilePath: 'Target File Path',
                            lowerXRangeIndex: 'Lower X Range Index',
                            upperXRangeIndex: 'Upper X Range Index',
                            lowerYRangeIndex: 'Lower Y Range Index',
                            upperYRangeIndex: 'Upper Y Range Index'
                        }
                    },
                    scheduling: {
                        header: 'Scheduling',
                        add: 'Add Schedule',
                        table: {
                            headers: {
                                id: 'ID', 
                                at: 'At',
                                cronExpression: 'Cron Expression',
                            },
                        }
                    }
                },
                steps: {
                    checkoutRepo: 'Checkout Repository',
                    shellScript: 'Shell Script',
                    pythonScript: 'Python Script',
                    scikitTrainingDataset: 'Load Training Dataset',
                    scikitTestingDataset: 'Load Testing Dataset',
                    scikitUnsupervisedModels: 'Scikit Unsupervised Models',
                    leastSquares: 'Least Squares',
                    ridgeRegression: 'Ridge Regression',
                    ridgeRegressionCV: 'Ridge Regression CV',
                    ridgeClassifier: 'Ridge Classifier',
                    ridgeClassifierCV: 'Ridge Classifier CV',
                    lasso: 'Lasso',
                    lassoCV: 'Lasso CV',
                    lassoLars: 'Lasso Lars',
                    lassoLarsCV: 'Lasso Lars CV',
                    lassoLarsIC: 'Lasso Lars IC',
                    multiTaskLasso: 'Multi Task Lasso',
                    multiTaskLassoCV: 'Multi Task Lasso CV',
                    elasticNet: 'Elastic Net',
                    elasticNetCV: 'Elastic Net CV',
                    multiTaskElasticNet: 'Multi Task Elastic Net',
                    multiTaskElasticNetCV: 'Multi Task Elastic Net CV',
                    lars: 'Lars',
                    larsCV: 'Lars CV',
                    omp: 'OMP',
                    ompCV: 'OMP CV',
                    bayesianRidge: 'Bayesian Ridge',
                    bayesianARD: 'Bayesian ARD',
                    logisticRegression: 'Logistic Regression',
                    logisticRegressionCV: 'Logistic Regression CV',
                    tweedieRegressor: 'Tweedie Regressor',
                    poissonRegressor: 'Poisson Regressor',
                    gammaRegressor: 'Gamma Regressor',
                    sgdClassifier: 'SGD Classifier',
                    sgdRegressor: 'SGD Regressor',
                    perceptron: 'Perceptron',
                    passiveAgressiveClassifier: 'Passive Agressive Classifier',
                    passiveAgressiveRegressor: 'Passive Agressive Regressor',
                    huberRegression: 'Huber Regression',
                    ransacRegression: 'Ransac Regression',
                    theilSenRegression: 'Theil Sen Regression',
                    quantileRegression: 'Quantile Regressiom',
                }
            },
            runs: {
                name: "Runs",
                header: "Pipeline {id} Runs",
                pipelineRuns: {
                    name: "Runs",
                    header: "Runs",
                    table: {
                        headers: {
                            name: "Name",
                            lastRun: "Last Run",
                        }
                    },
                    dialog: {
                        header: 'Create Run for Pipeline {id}?'
                    }
                },
                results: {
                    header: '{pipelineName} - Run {runID}',
                    dialog: {
                        edit: {
                            header: 'Edit {name} Step',
                        },
                    },
                    log: {
                        header: 'Log',
                    }
                },
                table: {
                    headers: {
                        id: "ID",
                        status: "Status",
                        created: "Created",
                        lastRun: "Last Run",
                    },
                    dialog: {
                        execute: {
                            header: 'Execute run {id}?',
                            body: 'This will erase all previous data associated with the run.',
                        }
                    }
                },
                dialog: {
                    create: {
                        header: 'Create Run for Pipeline ${id}?'
                    }
                },
            },
            profile: {
                name: "Profile",
                header: "Profile",
                user: {
                    greeting: 'Howdy, <b>{username}</b>!'
                },
                form: {
                    name: {
                        label: "Name",
                        help: "Required. Your name",
                    },
                    currentPassword: {
                        label: "Current Password",
                        help: "Required. Your current password",
                    },
                    newPassword: {
                        label: "New Password",
                        help: "Required. Your new password",
                    },
                    confirmNewPassword: {
                        label: "Confirm new Password",
                        help: "Required. Repeat your new password",
                    },
                    success: {
                        usernameChanged: 'Username changed with sucesss!',
                        passwordChanged: 'Password changed with sucesss!',
                    },
                    errors: {
                        newPassword: {
                            notEqual: "The new and confirmation passwords are not equal!"
                        }
                    },
                },
            },
            login: {
                name: 'Login',
                username: {
                    name: 'Username',
                    placeholder: 'Username',
                },
                password: {
                    name: 'Password',
                    placeholder: 'Minimum of 8 characters'
                },
                submit: 'Submit',
                remember: 'Remember me',
                signup: 'Sign Up'
            },
            signup: {
                name: 'Sign Up',
                username: {
                    name: 'Username',
                    placeholder: 'Username',
                },
                password: {
                    name: 'Password',
                    placeholder: 'Minimum of 8 characters'
                },
                submit: 'Submit',
            }
        },
        tables: {
            page: 'Page {page} of {count}',
        },
        buttons: {
            confirm: 'Confirm',
            save: 'Save',
            submit: 'Submit',
            cancel: 'Cancel',
            back: 'Back',
            next: 'Next',
        },
        messages: {
            types: {
                success: 'Success',
                error: 'Error',
            },
        },
        global: {
            app: {
                name: 'DI'
            },
            logout: 'Log Out',
            untitled: 'Untitled',
        }
    }
};

export const i18n = createI18n({
    legacy: false, // you must set `false`, to use Composition API
    locale: 'en', // set locale
    fallbackLocale: 'en', // set fallback locale
    messages
})