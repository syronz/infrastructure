import { Injectable } from '@angular/core';
import { environment as env } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import * as jwt_decode from 'jwt-decode';


@Injectable({
	providedIn: 'root'
})
export class AuthService {

	constructor(
		private http: HttpClient
	) { }

	logout() {
		localStorage.removeItem('currentUser');
	}

	login(username: string, password: string) {
		return this.http.post<any>(`${env.apiUrl}/auth`, { username, password })
			.pipe(map(user => {
				if (user && user.token) {
					localStorage.setItem('currentUser', JSON.stringify(user));
				}
				return user;
			}));
	}


	getUserJWT() {
		if(localStorage.getItem('currentUser')) {
			let userJWT = JSON.parse(localStorage.getItem('currentUser'));
			if ( !('token' in userJWT)){
				return false
			}

			userJWT.exp = this.getTokenExpiration(userJWT.token);
			userJWT.name = this.getTokenName(userJWT.token);

			return userJWT
		}
		return false
	}

	getTokenExpiration(token: string): Date {
		const decoded = jwt_decode(token);

		if (decoded.exp === undefined) return null;

		// For converting UNIX date to javascript date we have to multilpe it to 1000
		const date = new Date(decoded.exp * 1000);
		return date;
	}

	public getTokenName(token: string): string {
		const decoded = jwt_decode(token);
		if (decoded.name === undefined) return null;
		return decoded.name;
	}

	getProfile() {
		return this.http.get<any>(`${env.apiUrl}/auth/profile`);
	}

	updateProfile(userData: any) {
		return this.http.put<any>(`${env.apiUrl}/auth/profile`, userData.value);
	}
		

}




