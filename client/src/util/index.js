import axios from 'axios';
import moment from 'moment';

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
  return moment(date).format('MMMM Do YYYY, h:mm:ss a');
}

export const golangType = function(type) {
  if (type == 'text') {
      return 'string';
  }
  if (type == 'number') {
      return 'float64';
  }
  if (type == 'select') {
      return 'string';
  }
  if (type == 'text') {
      return 'string';
  }
  if (type == 'checkbox') {
    return 'bool';
  }
          
  return type
}