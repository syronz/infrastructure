import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DictionaryService } from '../../services/dictionary.service';
import { MatSnackBar } from '@angular/material';
import { environment as env } from '../../../environments/environment';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
	styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
	id: string;
	isNew: boolean;
	userInfo: any;
	userForm: FormGroup;
	error: any;

	constructor(
		private route: ActivatedRoute,
		private authService: AuthService,
		public loading: LoadingSpinnerService,
		private formBuilder: FormBuilder,
		private router: Router,
		private dict: DictionaryService,
		private snackBar: MatSnackBar
	) { }

	ngOnInit() {
		this.userForm = this.formBuilder.group({
			name: [''],
			language: [''],
			password: [''],
			confirmPassword: [''],
		});

		this.get();
	}

	get() {
		this.loading.show();
		this.authService.getProfile()
			.subscribe(
				response => {
					this.userInfo = response.data;
					this.userForm.controls['name'].setValue(this.userInfo.name);
					this.userForm.controls['language'].setValue(this.userInfo.language);
					this.loading.hide();
				},
				error => {
					console.warn('APPLICATION ERROR', error);
					this.loading.hide();
				}
			);
	}

	get f() {
		return this.userForm.controls;
	}

	onSubmit() {
		this.userForm.controls.name.setErrors(null);
		this.userForm.controls.password.setErrors(null);
		this.userForm.controls.confirmPassword.setErrors(null);
		this.userForm.controls.language.setErrors(null);
		this.error = {};
		if (this.userForm.invalid) {
			return;
		}

		if (this.f.password.value !== this.f.confirmPassword.value){
			this.f.password.setErrors({'incorrect': true});
			this.f.confirmPassword.setErrors({'incorrect': true});
			this.error = { error: { message: this.dict.translate("Passwords not match")} };
			return;
		}


		this.loading.show();
			this.authService.updateProfile(this.userForm)
				.subscribe(
					data => {
						//this.router.navigate(['/users']);
						this.snackBar.open( data.message, this.dict.translate("OK"), {
							duration: env.snackBarDuration,
						});
						this.loading.hide();
					},
					error => {
						this.loading.hide();
						console.warn('APPLICATION ERROR', error);
						this.error = error;
						let err = error.error.error;
						if (err.fields != null) {
							if ( err.fields.includes('name') ) {
								this.userForm.controls.name.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('password') ) {
								this.userForm.controls.password.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('language') ) {
								this.userForm.controls.language.setErrors({'incorrect': true});
							}
						}
					}
				);
	}


}
