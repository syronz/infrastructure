import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router, CanActivate, CanActivateChild } from '@angular/router';
import { Observable } from 'rxjs';
import { AuthService } from '../services/auth.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate, CanActivateChild  {
	constructor(
		private router: Router,
		private authService: AuthService,
	){}

	canActivate( next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
		if (!this.isLogged()){
			this.router.navigate(['/login']);
			return false;
		}

		return true;
	}

	canActivateChild( next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
		if (!this.isLogged()){
			this.router.navigate(['/login']);
			return false;
		}

		return true;
	}

	isLogged(): boolean {
		const userJWT = this.authService.getUserJWT()
		if (userJWT === false) return false;
		const now = new Date();
		if (userJWT.exp < now) return false;

		return true;
	}

}
