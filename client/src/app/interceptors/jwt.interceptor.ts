import { Injectable } from '@angular/core';
import { HttpRequest,HttpResponse, HttpErrorResponse, HttpHandler, HttpEvent, HttpInterceptor } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material';
import { environment as env } from '../../environments/environment';

@Injectable()
export class JwtInterceptor implements HttpInterceptor {

	constructor(
		private router: Router,
		private snackBar: MatSnackBar,
	){}

	

	intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
		let currentUser = JSON.parse(localStorage.getItem('currentUser'));
		if (currentUser && currentUser.token) {
			request = request.clone({
				setHeaders: {
					Authorization: `Bearer ${currentUser.token}`
				}
			});
		}

		return next.handle(request).pipe(
			map((event: HttpEvent<any>) => {
				//if (event instanceof HttpResponse) {
				//console.log('event--->>>', event);
				//}
				return event;
			}),
			catchError((error: HttpErrorResponse) => {
				if(error.status == 401){
					this.router.navigate(['login']);
				}
				if(error.status == 405) {
					this.snackBar.open(error.error.message, "X", {
						duration: env.snackBarDuration
					});
				}
				//return error;
				return throwError(error);
			}) );
	}
}


