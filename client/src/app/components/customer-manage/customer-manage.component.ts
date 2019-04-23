import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { CustomerService } from '../../services/customer.service';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DictionaryService } from '../../services/dictionary.service';
import { MatSnackBar } from '@angular/material';
import { environment as env } from '../../../environments/environment';

@Component({
  selector: 'app-customer-manage',
  templateUrl: './customer-manage.component.html',
  styleUrls: ['./customer-manage.component.css']
})
export class CustomerManageComponent implements OnInit {
	id: string;
	isNew: boolean;
	dataInfo: any;
	dataForm: FormGroup;
	error: any;

	constructor(
		private route: ActivatedRoute,
		private customerService: CustomerService,
		public loading: LoadingSpinnerService,
		private formBuilder: FormBuilder,
		private router: Router,
		private dict: DictionaryService,
		private snackBar: MatSnackBar
	) { }

	ngOnInit() {
		this.dataForm = this.formBuilder.group({
			title: [''],
			name: [''],
			phone1: [''],
			phone2: [''],
			createdAt: [''],
			detail: [''],
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
		this.customerService.getById(id)
			.subscribe(
				response => {
					this.dataInfo = response.data;
					this.dataForm.controls['title'].setValue(this.dataInfo.title);
					this.dataForm.controls['name'].setValue(this.dataInfo.name);
					this.dataForm.controls['phone1'].setValue(this.dataInfo.phone1);
					this.dataForm.controls['phone2'].setValue(this.dataInfo.phone2);
					this.dataForm.controls['createdAt'].setValue(this.dataInfo.created_at);
					this.dataForm.controls['detail'].setValue(this.dataInfo.detail);
					this.loading.hide();
				},
				error => {
					console.warn('APPLICATION ERROR', error);
					this.loading.hide();
				}
			);
	}

	get f() {
		return this.dataForm.controls;
	}

	onSubmit() {
		this.dataForm.controls.title.setErrors(null);
		this.dataForm.controls.name.setErrors(null);
		this.dataForm.controls.phone1.setErrors(null);
		this.dataForm.controls.phone2.setErrors(null);
		this.dataForm.controls.detail.setErrors(null);
		this.error = {};
		if (this.dataForm.invalid) {
			return;
		}

		this.loading.show();
		if(!this.isNew){
			this.customerService.update(this.dataForm, this.id)
				.subscribe(
					data => {
						this.router.navigate(['/customers']);
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
							if ( err.fields.includes('title') ) {
								this.dataForm.controls.title.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('detail') ) {
								this.dataForm.controls.detail.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('name') ) {
								this.dataForm.controls.name.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('phone1') ) {
								this.dataForm.controls.phone1.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('phone2') ) {
								this.dataForm.controls.phone2.setErrors({'incorrect': true});
							}
						}
					}
				);
		} else {
			this.customerService.create(this.dataForm)
				.subscribe(
					data => {
						this.router.navigate(['/customers']);
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
							if ( err.fields.includes('title') ) {
								this.dataForm.controls.title.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('detail') ) {
								this.dataForm.controls.detail.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('name') ) {
								this.dataForm.controls.name.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('phone1') ) {
								this.dataForm.controls.phone1.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('phone2') ) {
								this.dataForm.controls.phone2.setErrors({'incorrect': true});
							}
						}
					}
				);
		}
	}

	delete() {
		this.loading.show();
		this.customerService.delete(this.id)
			.subscribe(
				data => {
					this.loading.hide();
					this.router.navigate(['/customers']);
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
