import Axios from 'axios';
import { API_URL } from './api-config';

async function personAllForUser(): Promise<any> {
    const result = await Axios.get(`${API_URL}/person/allforuser`, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

async function personAdd(newPerson: Object): Promise<any> {
    const result = await Axios.post(`${API_URL}/person/add`, newPerson, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

async function personView(personId: String): Promise<any> {
    const result = await Axios.get(`${API_URL}/person/view/id/${personId}`, {
        headers: {
            // Authentication goes here
        }
    });
    return result.data;
}

export { personAllForUser, personAdd, personView }
