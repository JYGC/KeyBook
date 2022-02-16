import Axios from 'axios';
import { API_URL } from './api-config';

async function deviceAllForUser(): Promise<any> {
    const result = await Axios.get(`${API_URL}/device`, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

async function deviceAdd(newDevice: Object): Promise<any> {
    const result = await Axios.post(`${API_URL}/device/add`, newDevice, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

async function deviceView(deviceId: String): Promise<any> {
    const result = await Axios.get(`${API_URL}/device/view/id/${deviceId}`, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

async function deviceSave(existingDevice: Object, existingPerson: Object): Promise<any> {
    const result = await Axios.post(`${API_URL}/device/save`, {
        "device": existingDevice,
        "person": existingPerson
    }, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

export { deviceAllForUser, deviceAdd, deviceView, deviceSave }
