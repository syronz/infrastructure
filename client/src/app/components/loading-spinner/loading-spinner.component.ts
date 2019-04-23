import { Component, OnInit } from '@angular/core';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';

@Component({
  selector: 'app-loading-spinner',
  templateUrl: './loading-spinner.component.html',
  styleUrls: ['./loading-spinner.component.css']
})
export class LoadingSpinnerComponent implements OnInit {

	spinnerVisible: boolean;

	constructor(private loadingSpinnerService: LoadingSpinnerService) {
		this.spinnerVisible = this.loadingSpinnerService.isSpinnerVisable;

		this.loadingSpinnerService.getMessage().subscribe(
			msg => {
				this.spinnerVisible = msg;
			}
		);
		

	
	}

	ngOnInit() {
	}


}
