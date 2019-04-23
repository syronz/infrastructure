import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { UserService } from '../../services/user.service';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DictionaryService } from '../../services/dictionary.service';
import { MatSnackBar } from '@angular/material';
import { environment as env } from '../../../environments/environment';

@Component({
	selector: 'app-user-manage',
	templateUrl: './user-manage.component.html',
	styleUrls: ['./user-manage.component.css']
})
export class UserManageComponent implements OnInit {
	id: string;
	isNew: boolean;
	userInfo: any;
	userForm: FormGroup;
	error: any;

	constructor(
		private route: ActivatedRoute,
		private userService: UserService,
		public loading: LoadingSpinnerService,
		private formBuilder: FormBuilder,
		private router: Router,
		private dict: DictionaryService,
		private snackBar: MatSnackBar
	) { }

	ngOnInit() {
		this.userForm = this.formBuilder.group({
			name: [''],
			username: [''],
			role: [''],
			city: [''],
			director: [''],
			language: [''],
			password: [''],
			confirmPassword: [''],
		});

		this.id = this.route.snapshot.paramMap.get("id");
		if (this.id === "new"){
			this.isNew = true;
		}

		if (!this.isNew){
			this.get(this.id);
		}
	}

	get(id: string) {
		this.loading.show();
		this.userService.getById(id)
			.subscribe(
				response => {
					this.userInfo = response.data;
					this.userForm.controls['name'].setValue(this.userInfo.name);
					this.userForm.controls['username'].setValue(this.userInfo.username);
					this.userForm.controls['role'].setValue(this.userInfo.role);
					this.userForm.controls['city'].setValue(this.userInfo.city);
					this.userForm.controls['director'].setValue(this.userInfo.director);
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
		this.userForm.controls.username.setErrors(null);
		this.userForm.controls.name.setErrors(null);
		this.userForm.controls.password.setErrors(null);
		this.userForm.controls.role.setErrors(null);
		this.userForm.controls.city.setErrors(null);
		this.userForm.controls.director.setErrors(null);
		this.userForm.controls.language.setErrors(null);
		this.error = {};
		if (this.userForm.invalid) {
			return;
		}

		if (this.f.password.value !== this.f.confirmPassword.value){
			this.f.password.setErrors({'incorrect': true});
			//this.f.confirmPassword.setErrors({'incorrect': true});
			this.error = { error: { message: this.dict.translate("Passwords not match")} };
			return;
		}


		this.loading.show();
		if(!this.isNew){
			this.userService.update(this.userForm, this.id)
				.subscribe(
					data => {
						this.router.navigate(['/users']);
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
							if ( err.fields.includes('username') ) {
								this.userForm.controls.username.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('name') ) {
								this.userForm.controls.name.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('password') ) {
								this.userForm.controls.password.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('role') ) {
								this.userForm.controls.role.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('city') ) {
								this.userForm.controls.city.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('director') ) {
								this.userForm.controls.director.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('language') ) {
								this.userForm.controls.language.setErrors({'incorrect': true});
							}
						}
					}
				);
		} else {
			this.userService.create(this.userForm)
				.subscribe(
					data => {
						this.router.navigate(['/users']);
						this.snackBar.open(data.message, this.dict.translate("OK"), {
							duration: env.snackBarDuration,
						});
						this.loading.hide();
					},
					error => {
						this.loading.hide();
						console.warn('APPLICATION ERROR', error);
						this.error = error;
						let err = error.error.error;
						if( err.fields != null) {
							if (err.fields.includes('username')) {
								this.userForm.controls.username.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('name') ) {
								this.userForm.controls.name.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('password') ) {
								this.userForm.controls.password.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('role') ) {
								this.userForm.controls.role.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('city') ) {
								this.userForm.controls.city.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('director') ) {
								this.userForm.controls.director.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('language') ) {
								this.userForm.controls.language.setErrors({'incorrect': true});
							}
						}
					}
				);
		}
	}

	delete() {
		this.loading.show();
		this.userService.delete(this.id)
			.subscribe(
				data => {
					this.loading.hide();
					this.router.navigate(['/users']);
					this.snackBar.open(data.message, this.dict.translate("OK"), {
						duration: env.snackBarDuration,
					});
				},
				error => {
					this.loading.hide();
					console.warn("APPLICATION ERROR", error);
					this.snackBar.open(error.error.message, this.dict.translate("OK"), {
						duration: env.snackBarDuration,
					});
				}
			);
	}

}

