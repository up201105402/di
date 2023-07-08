import {
  mdiChartTimelineVariant,
  mdiRun,
  mdiAccount
} from "@mdi/js";

export default [
  {
    to: "/pipelines",
    icon: mdiChartTimelineVariant,
    label: "Pipelines",
  },
  {
    to: "/runs",
    icon: mdiRun,
    label: "Runs",
  },
  {
    to: "/profile",
    icon: mdiAccount,
    label: "Profile",
  },
  /*
  {
    label: "Dropdown",
    icon: mdiViewList,
    menu: [
      {
        label: "Item One",
      },
      {
        label: "Item Two",
      },
    ],
  },
  {
    href: "https://github.com/justboil/admin-one-vue-tailwind",
    label: "GitHub",
    icon: mdiGithub,
    target: "_blank",
  },
  {
    href: "https://github.com/justboil/admin-one-react-tailwind",
    label: "React version",
    icon: mdiReact,
    target: "_blank",
  },
  */
];
