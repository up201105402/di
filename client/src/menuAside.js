import {
  mdiChartTimelineVariant,
  mdiRun,
  mdiDataMatrix,
  mdiCodeBraces,
  mdiAccount
} from "@mdi/js";
import { i18n } from '@/i18n';

const { t } = i18n.global;

export default [
  {
    to: "/pipelines",
    icon: mdiChartTimelineVariant,
    label: t('pages.pipelines.name'),
  },
  {
    to: "/runs",
    icon: mdiRun,
    label: t('pages.runs.name'),
  },
  {
    to: "/datasets",
    icon: mdiDataMatrix,
    label: t('pages.datasets.name'),
  },
  {
    to: "/trainers",
    icon: mdiCodeBraces,
    label: t('pages.trainers.name'),
  },
  {
    to: "/testers",
    icon: mdiCodeBraces,
    label: t('pages.testers.name'),
  },
  {
    to: "/trained",
    icon: mdiCodeBraces,
    label: t('pages.trained.name'),
  },
  /*
  {
    to: "/profile",
    icon: mdiAccount,
    label: t('pages.profile.name'),
  },
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
