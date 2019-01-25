import {HttpService} from './helpers/HttpService.js';

var $ = document.querySelector.bind(document); 

var http = new HttpService();

const resultDiv = $('#resultDiv')

document.querySelector('#btnCat').onclick = (event) => {
  event.preventDefault();
  console.log("Clicou no Gato");
  let ret = http.get('${BACKEND_URL}');
  ret.then( ret => {
    var img = document.createElement("img");
    img.src = "cat.jpg"
    resultDiv.appendChild(img);
    //  resultDiv.innerHTML += "Resposta do servidor " + ret.animal + "<br />";
  }).catch(error => {
    console.log("Deu ruim: " + error);
    var img = document.createElement("img");
    img.src = "hedhog.jpg"
    resultDiv.appendChild(img);
  }) 
}

document.querySelector('#btnDog').onclick = (event) => {
  event.preventDefault();
  console.log("Clicou no Cachorro");
  let ret = http.get('${BACKEND_URL}');
  ret.then( ret => {
    var img = document.createElement("img");
    img.src = "dog.jpg";
    resultDiv.appendChild(img);
//    resultDiv.innerHTML += "Resposta do servidor " + ret.animal + "<br />";
  }).catch(error => { 
    console.log("Deu ruim: " + error)
    var img = document.createElement("img");
    img.src = "hedhog.jpg"
    resultDiv.appendChild(img);
  });
}

