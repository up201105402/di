import { createRouter, createWebHashHistory } from "vue-router";
import Home from "@/views/HomeView.vue";
import Pipelines from "@/views/PipelinesView.vue";
import Runs from "@/views/RunsView.vue";
import PipelineRuns from "@/views/PipelineRunsView.vue";
import PipelineEditor from "@/views/PipelineEditorView.vue";
import RunResults from "@/views/RunResultsView.vue";
import Feedback from "@/views/FeedbackView.vue";
import SingleFeedback from "@/views/SingleQueryFeedbackView.vue";
import Datasets from "@/views/DatasetsView.vue";

const routes = [
  // {
  //   meta: {
  //     title: "Select style",
  //   },
  //   path: "/",
  //   name: "style",
  //   component: Style,
  // },
  {
    meta: {
      title: "Dashboard",
    },
    path: "/dashboard",
    name: "dashboard",
    component: Home,
    private: true,
  },
  {
    meta: {
      title: "Pipelines",
    },
    path: "/pipelines",
    name: "pipelines",
    component: Pipelines,
    private: true,
  },
  {
    meta: {
      title: "Runs",
    },
    path: "/runs",
    name: "runs",
    component: Runs,
    private: true,
  },
  {
    meta: {
      title: "Pipeline Editor",
    },
    path: "/pipelines/edit/:id",
    name: "pipeline",
    component: PipelineEditor,
    private: true,
  },
  {
    meta: {
      title: "Pipeline Runs",
    },
    path: "/pipelines/runs/:id",
    name: "pipelineruns",
    component: PipelineRuns,
    private: true,
  },
  {
    meta: {
      title: "Run Results",
    },
    path: "/runresults/:id",
    name: "runsresults",
    component: RunResults,
    private: true,
  },
  {
    meta: {
      title: "Human Feedback",
    },
    path: "/feedback/:id",
    name: "feedback",
    component: Feedback,
    private: true,
  },
  {
    meta: {
      title: "Single Query Human Feedback",
    },
    path: "/feedback/:id/query/:queryId",
    name: "singleFeedback",
    component: SingleFeedback,
    private: true,
  },
  {
    meta: {
      title: "Datasets",
    },
    path: "/datasets",
    name: "datasets",
    component: Datasets,
    private: true,
  },
  {
    meta: {
      title: "Tables",
    },
    path: "/tables",
    name: "tables",
    component: () => import("@/views/TablesView.vue"),
  },
  {
    meta: {
      title: "Forms",
    },
    path: "/forms",
    name: "forms",
    component: () => import("@/views/FormsView.vue"),
  },
  {
    meta: {
      title: "Profile",
    },
    path: "/profile",
    name: "profile",
    component: () => import("@/views/ProfileView.vue"),
    private: true,
  },
  {
    meta: {
      title: "Ui",
    },
    path: "/ui",
    name: "ui",
    component: () => import("@/views/UiView.vue"),
  },
  {
    meta: {
      title: "Responsive layout",
    },
    path: "/responsive",
    name: "responsive",
    component: () => import("@/views/ResponsiveView.vue"),
  },
  {
    meta: {
      title: "Login",
    },
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginView.vue"),
  },
  {
    meta: {
      title: "Sign Up",
    },
    path: "/signup",
    name: "signup",
    component: () => import("@/views/SignupView.vue"),
  },
  {
    meta: {
      title: "Error",
    },
    path: "/error",
    name: "error",
    component: () => import("@/views/ErrorView.vue"),
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/pipelines'
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 };
  },
});

export default router;
