import { Component, OnInit, ViewChild } from '@angular/core';
import { CityService } from '../../services/city.service';
import { PageEvent } from '@angular/material';
import { MatSort } from '@angular/material';
import { environment as env } from '../../../environments/environment';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';

@Component({
  selector: 'app-cities',
  templateUrl: './cities.component.html',
  styleUrls: ['./cities.component.css']
})
export class CitiesComponent implements OnInit {
	@ViewChild('searchBar') searchBar;
	@ViewChild('paginator') paginator;
	count: number;
	perPage: number;
	perPageOptions: number[] = env.perPageOptions;
	searchStr: string = "";
	sortField: string = null;
	sortDirection: string = null;
	pageEvent;


	displayedColumns: string[] = ['id', 'governorate', 'city', 'view'];
	dataSource: any;


	constructor(
		private cityService: CityService,
		public loading: LoadingSpinnerService,
	) { 
		this.perPage = this.perPageOptions[0];
	}

	ngOnInit() {
		this.listData(this.perPage,this.paginator.pageIndex + 1);
	}

	listData(perPage, page) {
		let searchStr = this.searchBar.nativeElement.value;
		this.loading.show()
		this.cityService.list(perPage, page, searchStr, this.sortField, this.sortDirection)
			.subscribe(
				response => {
					this.count = response.count;
					this.dataSource = (response.data);
					this.loading.hide();
				},
				error => {
					console.warn('APPLICATION ERROR', error);
					this.loading.hide();
				}
			); 
	}

	doSearch(){
		this.paginator.pageIndex = 0;
		this.listData(this.perPage, this.paginator.pageIndex);
	}

	loadPaginator(event?:PageEvent){
		this.listData(event.pageSize, event.pageIndex + 1);
		return event;
	}

	sortData(event) {
		this.paginator.pageIndex = 0;
		this.sortField = event.active;
		this.sortDirection = event.direction;
		this.listData(this.paginator.pageSize, 1);
	}

}
