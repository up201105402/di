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
