import {toast} from "vue3-toastify";

const notify = (message: string, success = true) => {
    if (success) {
        toast.success(message);
    }
    else {
        toast.error(message);
    }
};

export {notify};