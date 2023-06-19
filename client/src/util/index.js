import axios from 'axios';

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
  .trim()