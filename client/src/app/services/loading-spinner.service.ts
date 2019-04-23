import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';


@Injectable({
  providedIn: 'root'
})
export class LoadingSpinnerService {

	isSpinnerVisable: boolean;

	spinnerVisibilityChange: Subject<boolean> = new Subject<boolean>();

	constructor()  {
		this.spinnerVisibilityChange.subscribe((value) => {
			this.isSpinnerVisable = value
		});
	}

	toggleSidebarVisibility() {
		this.spinnerVisibilityChange.next(!this.isSpinnerVisable);
	}

	show() {
		this.spinnerVisibilityChange.next(true);
	}

	hide() {
		this.spinnerVisibilityChange.next(false);
	}

	getMessage(): Observable<any> {
		return this.spinnerVisibilityChange.asObservable();
	}




}
