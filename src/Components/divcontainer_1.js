import React from 'react'
import './divcontainer_1.css'
import ImgAsset from '../public'
export default function Divcontainer_1 () {
	return (
		<div className='divcontainer_1_divcontainer'>
			<div className='divflexframe'>
				<div className='divsidecolumn'>
					<div className='divbigcard'>
						<img className='_64a6cdff000962bbfb4a9cea_card1p800jpg' src = {ImgAsset.divcontainer_1__64a6cdff000962bbfb4a9cea_card1p800jpg} />
					</div>
				</div>
				<div className='divmiddlecolumn'>
					<span className='AddYours'><br/>Add Yours </span>
					<div className='Link'>
						<div className='divstickycontent'>
							<div className='divroundbtnborder'/>
							<div className='divgradientbtninner'>
								<span className='Add'>Add</span>
							</div>
						</div>
					</div>
				</div>
				<div className='divsidecolumn_1'>
					<div className='divbigcard_1'>
						<img className='_64a6cdff000962bbfb4a9d10_stewartmacleanZs1WKNa4Oy0unsplashp800jpg' src = {ImgAsset.divcontainer_1__64a6cdff000962bbfb4a9d10_stewartmacleanZs1WKNa4Oy0unsplashp800jpg} />
					</div>
					<span className='Andreceiveabloggerrewardpayment'>And receive a blogger reward payment</span>
				</div>
			</div>
		</div>
	)
}