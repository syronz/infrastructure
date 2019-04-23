import { Injectable } from '@angular/core';
import { environment as env } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ActivityService {

	constructor(
		private http: HttpClient
	) { }

	list(perPage, page, searchStr, sortField, sortDirection) {
		var params: any = {
			page: page,
			perPage: perPage,
		};

		if(searchStr) params.search = searchStr;
		if(sortDirection) {
			params.sortField = sortField;
			params.sortDirection = sortDirection;
		}

		return this.http.get<any>(`${env.apiUrl}/activities`, { params: params });

	}
	
}
