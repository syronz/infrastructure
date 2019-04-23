import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';
import { first } from 'rxjs/operators';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
	loginForm: FormGroup;
	loading = false;
	submitted = false;
	hide = true;
	returnUrl: string;
	error = '';


	constructor(
		private formBuilder: FormBuilder,
		private authService: AuthService,
		private router: Router,
	) { 
		this.authService.logout();
	}

  ngOnInit() {
		this.loginForm = this.formBuilder.group({
			username: ['', Validators.required],
			password: ['', Validators.required]
		});
  }

	get f() {
		return this.loginForm.controls;
	}

	onSubmit() {
		this.submitted = true;
		this.loginForm.controls.username.setErrors(null);
		this.loginForm.controls.password.setErrors(null);

		if (this.loginForm.invalid) {
			return ;
		}

		this.loading = true;

		this.authService.login(this.f.username.value, this.f.password.value)
			.pipe(first())
			.subscribe(
				data => {
					this.router.navigate(['/'])
				},
				error => {
					this.loginForm.controls.username.setErrors({'incorrect': true});
					this.loginForm.controls.password.setErrors({'incorrect': true});
					this.error = error;
					this.loading = false;
				});
			
	}

}
