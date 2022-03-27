import Axios from 'axios';

export class _ApiCall {
    protected _apiPath: string;

    private __successCallback: (response) => any;
    set successCallback(value: (response) => any) {
        this.__successCallback = value;
    }

    private __failedCallback: (any) => any;
    set failedCallback(value: (any) => any) {
        this.__failedCallback = value;
    }

    call() {
        Axios.get(this._apiPath)
          .then(response => this.__successCallback(response))
          .catch(e => this.__failedCallback(e));
    }
}