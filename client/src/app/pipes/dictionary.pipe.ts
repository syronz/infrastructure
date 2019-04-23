import { Pipe, PipeTransform } from '@angular/core';
import { DictionaryService } from '../services/dictionary.service';

@Pipe({
	name: 't',
	//pure: false
})
export class DictionaryPipe implements PipeTransform {

	constructor(private dict: DictionaryService){}

  transform(value: string): string {
		return this.dict.translate(value);
  }


}
