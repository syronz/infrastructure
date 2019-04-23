import { Injectable } from '@angular/core';
import * as Sentences from '../../assets/dictionary/sentences.json';


@Injectable({
  providedIn: 'root'
})
export class DictionaryService {
	sentences: any = Sentences;
	currentLanguage: string = "ku";

	constructor() {
		this.loadWords();
	}

	loadWords() {
		let currentUser = JSON.parse(localStorage.getItem('currentUser'));
		if('language' in currentUser){
			if(['en', 'ku', 'ar'].includes(currentUser.language)){
				this.currentLanguage = currentUser.language;
			}
		}
	}

	public translate(str: string): string {
		if ( !(str in this.sentences.default) ){
			console.error(`NEED TRANSLATION: ${str}`);
			return `※※※ ${str} ※※※`;
		}

		return this.sentences.default[str][this.currentLanguage];
	}


}
