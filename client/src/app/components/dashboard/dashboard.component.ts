import { Component, OnInit, ChangeDetectorRef, ViewChild } from '@angular/core';
import { MediaMatcher } from '@angular/cdk/layout';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
	@ViewChild('snav') snav;
	mobileQuery: MediaQueryList;
	currentSideStatus: boolean = true;


	private _mobileQueryListener: () => void;

	constructor(
		changeDetectorRef: ChangeDetectorRef,
		media: MediaMatcher,
		
	) {
		this.mobileQuery = media.matchMedia('(max-width: 600px)');
		this._mobileQueryListener = () => changeDetectorRef.detectChanges();
		this.mobileQuery.addListener(this._mobileQueryListener);
	}

  ngOnInit() {
  }

	showMSG = function(v) {
		this.snav.toggle();
	}

}
