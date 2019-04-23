import { Injectable } from '@angular/core';
import { environment as env } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class DirectorService {
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

		return this.http.get<any>(`${env.apiUrl}/directors`, { params: params });
	}

	getById(id: string) {
		return this.http.get<any>(`${env.apiUrl}/directors/${id}`);
	}

	update(data: any, id: string) {
		return this.http.put<any>(`${env.apiUrl}/directors/${id}`, data.value);
	}

	create(data: any) {
		return this.http.post<any>(`${env.apiUrl}/directors`, data.value);
	}

	delete(id: string) {
		return this.http.delete<any>(`${env.apiUrl}/directors/${id}`);
	}

}
