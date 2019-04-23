import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { DirectorService } from '../../services/director.service';
import { LoadingSpinnerService } from '../../services/loading-spinner.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DictionaryService } from '../../services/dictionary.service';
import { MatSnackBar } from '@angular/material';
import { environment as env } from '../../../environments/environment';

@Component({
  selector: 'app-director-manage',
  templateUrl: './director-manage.component.html',
  styleUrls: ['./director-manage.component.css']
})
export class DirectorManageComponent implements OnInit {
	id: string;
	isNew: boolean;
	dataInfo: any;
	dataForm: FormGroup;
	error: any;

	constructor(
		private route: ActivatedRoute,
		private directorService: DirectorService,
		public loading: LoadingSpinnerService,
		private formBuilder: FormBuilder,
		private router: Router,
		private dict: DictionaryService,
		private snackBar: MatSnackBar
	) { }

	ngOnInit() {
		this.dataForm = this.formBuilder.group({
			detail: [''],
			director: [''],
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
		this.directorService.getById(id)
			.subscribe(
				response => {
					this.dataInfo = response.data;
					this.dataForm.controls['detail'].setValue(this.dataInfo.detail);
					this.dataForm.controls['director'].setValue(this.dataInfo.director);
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
		this.dataForm.controls.detail.setErrors(null);
		this.dataForm.controls.director.setErrors(null);
		this.error = {};
		if (this.dataForm.invalid) {
			return;
		}

		this.loading.show();
		if(!this.isNew){
			this.directorService.update(this.dataForm, this.id)
				.subscribe(
					data => {
						this.router.navigate(['/directors']);
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
							if ( err.fields.includes('detail') ) {
								this.dataForm.controls.detail.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('director') ) {
								this.dataForm.controls.director.setErrors({'incorrect': true});
							}
						}
					}
				);
		} else {
			this.directorService.create(this.dataForm)
				.subscribe(
					data => {
						this.router.navigate(['/directors']);
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
							if ( err.fields.includes('detail') ) {
								this.dataForm.controls.detail.setErrors({'incorrect': true});
							}
							if ( err.fields.includes('director') ) {
								this.dataForm.controls.director.setErrors({'incorrect': true});
							}
						}
					}
				);
		}
	}

	delete() {
		this.loading.show();
		this.directorService.delete(this.id)
			.subscribe(
				data => {
					this.loading.hide();
					this.router.navigate(['/directors']);
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
