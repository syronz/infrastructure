import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { CityService } from '../../services/city.service';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DictionaryService } from '../../services/dictionary.service';
import { MatSnackBar } from '@angular/material';
import { environment as env } from '../../../environments/environment';

@Component({
	selector: 'app-city-manage',
	templateUrl: './city-manage.component.html',
	styleUrls: ['./city-manage.component.css']
})
export class CityManageComponent implements OnInit {
	id: string;
	isNew: boolean;
	dataInfo: any;
	dataForm: FormGroup;
	error: any;

	constructor(
		private route: ActivatedRoute,
		private cityService: CityService,
		public loading: LoadingSpinnerService,
		private formBuilder: FormBuilder,
		private router: Router,
		private dict: DictionaryService,
		private snackBar: MatSnackBar
	) { }

	ngOnInit() {
		this.dataForm = this.formBuilder.group({
			governorate: [''],
			city: [''],
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
		this.cityService.getById(id)
			.subscribe(
				response => {
					this.dataInfo = response.data;
					this.dataForm.controls['governorate'].setValue(this.dataInfo.governorate);
					this.dataForm.controls['city'].setValue(this.dataInfo.city);
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
		this.dataForm.controls.governorate.setErrors(null);
		this.dataForm.controls.city.setErrors(null);
		this.error = {};
		if (this.dataForm.invalid) {
			return;
		}

		this.loading.show();
		if(!this.isNew){
			this.cityService.update(this.dataForm, this.id)
				.subscribe(
					data => {
						this.router.navigate(['/cities']);
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
							if ( err.fields.includes('governorate') ) {
								this.dataForm.controls.governorate.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('city') ) {
								this.dataForm.controls.city.setErrors({'incorrect': true});
							}
						}
					}
				);
		} else {
			this.cityService.create(this.dataForm)
				.subscribe(
					data => {
						this.router.navigate(['/cities']);
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
							if ( err.fields.includes('governorate') ) {
								this.dataForm.controls.governorate.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('city') ) {
								this.dataForm.controls.city.setErrors({'incorrect': true});
							}
						}
					}
				);
		}
	}

	delete() {
		this.loading.show();
		this.cityService.delete(this.id)
			.subscribe(
				data => {
					this.loading.hide();
					this.router.navigate(['/cities']);
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
