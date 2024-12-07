import React from 'react'
import './divfooterborder.css'
import ImgAsset from '../public'
export default function Divfooterborder () {
	return (
		<div className='divfooterborder_divfooterborder'>
			<div className='divfooterbackground'>
				<div className='divglobalpadding'>
					<div className='divfootercontent'>
						<div className='divfooterlists'>
							<span className='Socails'>Socails</span>
							<div className='divfooterlinkscontainer'>
								<div className='Link'>
									<span className='Twitter'>Twitter</span>
								</div>
								<div className='Link_1'>
									<span className='Facebook'>Facebook</span>
								</div>
								<div className='Link_2'>
									<span className='Instagtam'>Instagtam</span>
								</div>
							</div>
						</div>
						<div className='divfooterlists_1'>
							<span className='Contact'>Contact</span>
							<div className='divfooterlinkscontainer_1'>
								<div className='Link_3'>
									<span className='Address'>Address</span>
								</div>
								<div className='Link_4'>
									<span className='Email'>Email</span>
								</div>
								<div className='Link_5'>
									<span className='xxxxxxxxx'>+xxxxxxxxx  </span>
								</div>
							</div>
						</div>
						<div className='divfooterlists_2'>
							<span className='Links'>Links</span>
							<div className='divfooterlinkscontainer_2'>
								<div className='Link_6'>
									<span className='License'>License</span>
								</div>
								<div className='Link_7'>
									<span className='YouTube'>YouTube</span>
								</div>
								<div className='Link_8'>
									<span className='TrustPilot'>TrustPilot</span>
								</div>
							</div>
						</div>
					</div>
					<div className='divfootercta'>
						<div className='Link_9'>
							<div className='divroundbtnborder'/>
							<div className='divgradientbtninner'>
								<span className='GetClixr'>Get Clixr</span>
							</div>
						</div>
					</div>
					<div className='divfootercontent_1'>
						<div className='divcreditscontainer'>
							<div className='Link_10'>
								<span className='DesignDevelopmentby'>Design & Development by </span>
								<div className='spanunderlinedgreen'>
									<span className='Deveb'>Deveb</span>
									<div className='_64a6cdff000962bbfb4a9cb8_underlinegreensvg'>
										<img className='Vector' src = {ImgAsset.divfooterborder_Vector} />
									</div>
								</div>
							</div>
							<div className='divtextxsmall'>
								<span className='PoweredbyWebflow'>Powered by Webflow</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	)
}