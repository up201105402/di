import axios from 'axios';

// doRequest is a helper function for
// handling axios responses - reqOptions follow axios req config
export const doRequest = async (reqOptions) => {
    let error;
    let data;

    try {
        const response = await axios.request(reqOptions);
        data = response.data;
    } catch (e) {
        if (e.response) {
            error = e.response.data.error;
        } else if (e.request) {
            error = e.request;
        } else {
            error = e;
        }
    }

    return {
        data,
        error,
    };
};
