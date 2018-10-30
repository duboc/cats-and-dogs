import {HttpService} from './helpers/HttpService.js';

var $ = document.querySelector.bind(document); 

var http = new HttpService();

const resultDiv = $('#resultDiv')

document.querySelector('#btnCat').onclick = (event) => {
  event.preventDefault();
  console.log("Clicou no Gato");
  let ret = http.get('http://localhost:9090/api/cat');
  ret.then( ret => {
      resultDiv.innerHTML += "Resposta do servidor " + ret.animal + "<br />";
  }).catch(error => console.log("Deu ruim: " + error));
}

document.querySelector('#btnDog').onclick = (event) => {
  event.preventDefault();
  console.log("Clicou no Cachorro");
  let ret = http.get('http://localhost:9090/api/dog');
  ret.then( ret => {
      resultDiv.innerHTML += "Resposta do servidor " + ret.animal + "<br />";
  }).catch(error => console.log("Deu ruim: " + error));
}

