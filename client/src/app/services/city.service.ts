import { Injectable } from '@angular/core';
import { environment as env } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class CityService {

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

		return this.http.get<any>(`${env.apiUrl}/cities`, { params: params });
	}

	getById(id: string) {
		return this.http.get<any>(`${env.apiUrl}/cities/${id}`);
	}

	update(data: any, id: string) {
		return this.http.put<any>(`${env.apiUrl}/cities/${id}`, data.value);
	}

	create(data: any) {
		return this.http.post<any>(`${env.apiUrl}/cities`, data.value);
	}

	delete(id: string) {
		return this.http.delete<any>(`${env.apiUrl}/cities/${id}`);
	}


}
