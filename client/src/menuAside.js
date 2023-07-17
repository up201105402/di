import {
  mdiChartTimelineVariant,
  mdiRun,
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
    to: "/profile",
    icon: mdiAccount,
    label: t('pages.profile.name'),
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
