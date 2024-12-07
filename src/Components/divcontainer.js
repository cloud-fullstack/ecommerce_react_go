import React from 'react'
import './divcontainer.css'
import ImgAsset from '../public'
export default function Divcontainer () {
	return (
		<div className='divcontainer_divcontainer'>
			<div className='divaligncenter'>
				<div className='divshadowtoptext'>
					<span className='BuyerBlogs'>Buyer. Blogs.</span>
				</div>
			</div>
			<div className='divlogolooping'>
				<span className='Createdbyverifiedusersofthisproduct'>Created by verified  users of this product.</span>
				<div className='Link'>
					<div className='divsmallfont'>
						<span className='Rezaverseiswiththebloggers'>Rezaverse is with the bloggers. </span>
						<div className='_64a6cdff000962bbfb4a9cb8_underlinegreensvg'>
							<img className='Vector' src = {ImgAsset.divcontainer_Vector} />
						</div>
					</div>
					<div className='_64a6cdff000962bbfb4a9cba_greenarrowsvg'>
					</div>
				</div>
			</div>
			<div className='divshowcaselooping'>
				<div className='divshowcaserow'>
					<div className='divcard'>
					</div>
					<div className='divcard_1'>
						<img className='_64a6cdff000962bbfb4a9ccb_sc1p500jpg' src = {ImgAsset.divcard__64a6cdff000962bbfb4a9ccb_sc1p500jpg} />
						<img className='Star1' src = {ImgAsset.divcontainer_Star1} />
					</div>
					<div className='divcard_2'>
						<img className='_64a6cdff000962bbfb4a9cd0_mainp500jpg' src = {ImgAsset.divcard_1__64a6cdff000962bbfb4a9cd0_mainp500jpg} />
						<img className='Star1_1' src = {ImgAsset.divcontainer_Star1_1} />
					</div>
					<div className='divcard_3'>
						<img className='_64a6cdff000962bbfb4a9cd5_sc3p500jpg' src = {ImgAsset.divcard_2__64a6cdff000962bbfb4a9cd5_sc3p500jpg} />
						<img className='Star1_2' src = {ImgAsset.divcontainer_Star1_2} />
					</div>
				</div>
				<span className='ByName'>By Name</span>
				<div className='divshowcaserow_1'>
					<div className='divcard_4'>
					</div>
					<div className='divcard_5'>
						<img className='_64a6cdff000962bbfb4a9ce2_sc6p500jpg' src = {ImgAsset.divcontainer__64a6cdff000962bbfb4a9ce2_sc6p500jpg} />
						<img className='Star1_3' src = {ImgAsset.divcontainer_Star1_3} />
					</div>
					<div className='divcard_6'>
						<img className='_64a6cdff000962bbfb4a9cf4_sc7p500jpg' src = {ImgAsset.divcontainer__64a6cdff000962bbfb4a9cf4_sc7p500jpg} />
						<img className='Star1_4' src = {ImgAsset.divcontainer_Star1_4} />
					</div>
					<div className='divcard_7'>
						<img className='_64a6cdff000962bbfb4a9ce6_sc8p500jpg' src = {ImgAsset.divcontainer__64a6cdff000962bbfb4a9ce6_sc8p500jpg} />
						<img className='Star1_5' src = {ImgAsset.divcontainer_Star1_5} />
					</div>
					<div className='divcard_8'>
						<img className='_64a6cdff000962bbfb4a9cef_sc5p500jpg' src = {ImgAsset.divcontainer__64a6cdff000962bbfb4a9cef_sc5p500jpg} />
						<img className='Star1_6' src = {ImgAsset.divcontainer_Star1_6} />
					</div>
					<img className='divgradienttop' src = {ImgAsset.divcontainer_divgradienttop} />
				</div>
				<span className='ByName_1'>By Name</span>
				<span className='ByName_2'>By Name</span>
			</div>
		</div>
	)
}