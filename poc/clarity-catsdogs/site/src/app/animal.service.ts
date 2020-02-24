import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import 'rxjs/add/operator/map';

@Injectable()
export class AnimalService {

  constructor(private http: Http) { }

  get(){
    let url = `http://cd.vmwlatam.com:9090/api/cat`;
    return this.http.get(url).map(data => data.json());
  }
}
