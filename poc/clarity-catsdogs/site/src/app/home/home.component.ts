import { Component } from "@angular/core";
import { AnimalService } from "app/animal.service";
import { State } from "@clr/angular";

@Component({
    styleUrls: ['./home.component.scss'],
    templateUrl: './home.component.html',
    providers: [AnimalService],
})

export class HomeComponent {
    
    constructor(private animalService: AnimalService) {
    }

    pet = [];
    total = 0;

    onClickCat() {
        this.animalService.get().subscribe(data => {
            this.pet = data.results; 
            this.total = data.count;
            alert(this.pet[0]);
            var myImg = "app/cat.jpg";
            if(this.pet[0]['animal'] == "cat"){
                var myImg = "app/cat.jpg";
                alert("gato");
            };
        });
    }

    onClickDog(){
        alert("Você clicou no cachorro");
    }
}
