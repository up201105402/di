import axios from 'axios';
import moment from 'moment';
import { i18n } from '@/i18n';

const { t } = i18n.global;

export const doRequest = async (reqOptions) => {
  let status;
  let error;
  let data;

  try {
    const response = await axios.request(reqOptions);
    data = response.data;
    status = response.status;
  } catch (e) {
    if (e.response) {
      error = e.response.data.error;
      status = e.response.status;
    } else if (e.request) {
      error = e.request;
      status = e.request.status;
    } else {
      error = e;
    }
  }

  return {
    status,
    data,
    error,
  };
};

export const cleanObject = (object) => {
  Object.keys(object).forEach(key => {
    if (object[key] == null) {
      delete object[key];
    }
  });
}

export { customDelay, camel2title }

const customDelay = () => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve()
    }, 500)
  })
}

const camel2title = (str) => str
  .replace(/([A-Z])/g, (match) => ` ${match}`)
  .replace(/^./, (match) => match.toUpperCase())
  .replace('C V', 'CV')
  .replace('I C', 'IC')
  .trim()

export const removeDuplicates = (arr, uniqueProp) => {
  const uniqueProps = [];

  return arr.filter(element => {
    const isDuplicate = uniqueProps.includes(element[uniqueProp]);

    if (!isDuplicate) {
      uniqueProps.push(element[uniqueProp]);

      return true;
    }

    return false;
  });

}

export const formatDate = (date) => {
  const m = moment(date);

  if (!date || m.year() == 0) {
    return "-"
  }
  
  return moment(date).format('MMMM Do YYYY, h:mm:ss a');
}

export const parsePipelineDefinition = (entity, toast) => {
  try {
    return JSON.parse(entity.definition)
  } catch (e) {
    if (entity.definition != "") {
      toast.add({ severity: 'error', summary: t('global.errors.parsing.header'), detail: t('global.errors.parsing.detail'), life: 3000 });
    }
    return [];
  }
}

export function validateCron(value) {
    if (!value) {
      return 'Cron expression is required!';
    }

    if (!value.match('(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\\d+(ns|us|µs|ms|s|m|h))+)|((((\\d+,)+\\d+|(\\d+(\\/|-)\\d+)|\\d+|\\*) ?){5,7})')) {
      return 'Cron expresion is not valid!';
    }

    return true;
};

export const deepFilterMenuBarSteps = (subject, key, value) => {
  if (subject instanceof Array) {
    return subject.map(elem => deepFilterMenuBarSteps(elem, key, value))
                  .filter(elem => elem != null)[0];
  }

  if (subject[key] == value) {
    return subject;
  }

  const items = subject['items'];

  if (items && items.length) {
    return items.map(item => deepFilterMenuBarSteps(item, key, value))
                .filter(elem => elem != null)[0];
  }

  return null
}

export const i18nFromStepName = (stepName) => {
  return t('pages.pipelines.edit.dialog.' + stepName + '.label');
}

export const getStatusTagSeverity = (statusID) => {
  switch (statusID) {
      case 1:
          return "info";
      case 2: 
          return "info";
      case 3: 
          return "danger";
      case 4:
          return "success";
      case 5:
          return "warning";
  }
  
  return "info";
}