import axios from 'axios';

// doRequest is a helper function for
// handling axios responses - reqOptions follow axios req config
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

export { customAxios, camel2title }

// This is just a mock of an actual axios instance.
const customAxios = {
  post: (formData) => {
    return new Promise((resolve, reject) => {
      let response = { status: 200 }
      if (formData.organizationInfo.org_name.toLowerCase().trim() !== 'formkit') {
        response = {
          status: 400,
          formErrors: ['There was an error in this form, please correct and re-submit to validate.'],
          fieldErrors: {
              'organizationInfo.org_name': 'Organization Name must be "FormKit", we tricked you!'
          }
        }
      }
      setTimeout(() => {
        if (response.status === 200) {
          resolve(response)
        } else {
          reject(response)
        }
      }, 1500)
    })
    
  }
}

const camel2title = (str) => str
  .replace(/([A-Z])/g, (match) => ` ${match}`)
  .replace(/^./, (match) => match.toUpperCase())
  .trim()