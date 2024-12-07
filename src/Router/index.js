import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import HomePage from '../Components/index';
import Frame from '../Components/Frame';
import Divmxauto from '../Components/divmxauto';
import Divmt8 from '../Components/divmt8';
import Divflex from '../Components/divflex';
import Iconamoonprofilecirclethin from '../Components/iconamoonprofilecirclethin';
import Iconamoonprofilecirclethin_1 from '../Components/iconamoonprofilecirclethin_1';
import Frame_1 from '../Components/Frame_1';
import Divhidden from '../Components/divhidden';
import Frame_2 from '../Components/Frame_2';
import Divmxauto_1 from '../Components/divmxauto_1';
import Divcontainer from '../Components/divcontainer';
import Divcontainer_1 from '../Components/divcontainer_1';
import Divcard from '../Components/divcard';
import Divcard_1 from '../Components/divcard_1';
import Divcard_2 from '../Components/divcard_2';
import Frame_3 from '../Components/Frame_3';
import Divcontainer_2 from '../Components/divcontainer_2';
import Divcontainer_3 from '../Components/divcontainer_3';
import Divwhitebackground from '../Components/divwhitebackground';
import DivFeatures from '../Components/divFeatures';
import Divglobalpadding from '../Components/divglobalpadding';
import DivServices from '../Components/divServices';
import DivFAQ from '../Components/divFAQ';
import Divglobalpadding_1 from '../Components/divglobalpadding_1';
import Divfooterborder from '../Components/divfooterborder';
import Divcard_3 from '../Components/divcard_3';
import Divcard_4 from '../Components/divcard_4';
import Divcard_5 from '../Components/divcard_5';
import Frame_4 from '../Components/Frame_4';
import Arrowright from '../Components/arrowright';
import Product_card from '../Components/product_card';
const RouterDOM = () => {
	return (
		<Router>
			<Switch>
				<Route exact path="/"><HomePage /></Route>
				<Route exact path="/frame"><Frame /></Route>
				<Route exact path="/divmxauto"><Divmxauto /></Route>
				<Route exact path="/divmt8"><Divmt8 /></Route>
				<Route exact path="/divflex"><Divflex /></Route>
				<Route exact path="/iconamoonprofilecirclethin"><Iconamoonprofilecirclethin /></Route>
				<Route exact path="/iconamoonprofilecirclethin_1"><Iconamoonprofilecirclethin_1 /></Route>
				<Route exact path="/frame_1"><Frame_1 /></Route>
				<Route exact path="/divhidden"><Divhidden /></Route>
				<Route exact path="/frame_2"><Frame_2 /></Route>
				<Route exact path="/divmxauto_1"><Divmxauto_1 /></Route>
				<Route exact path="/divcontainer"><Divcontainer /></Route>
				<Route exact path="/divcontainer_1"><Divcontainer_1 /></Route>
				<Route exact path="/divcard"><Divcard /></Route>
				<Route exact path="/divcard_1"><Divcard_1 /></Route>
				<Route exact path="/divcard_2"><Divcard_2 /></Route>
				<Route exact path="/frame_3"><Frame_3 /></Route>
				<Route exact path="/divcontainer_2"><Divcontainer_2 /></Route>
				<Route exact path="/divcontainer_3"><Divcontainer_3 /></Route>
				<Route exact path="/divwhitebackground"><Divwhitebackground /></Route>
				<Route exact path="/divfeatures"><DivFeatures /></Route>
				<Route exact path="/divglobalpadding"><Divglobalpadding /></Route>
				<Route exact path="/divservices"><DivServices /></Route>
				<Route exact path="/divfaq"><DivFAQ /></Route>
				<Route exact path="/divglobalpadding_1"><Divglobalpadding_1 /></Route>
				<Route exact path="/divfooterborder"><Divfooterborder /></Route>
				<Route exact path="/divcard_3"><Divcard_3 /></Route>
				<Route exact path="/divcard_4"><Divcard_4 /></Route>
				<Route exact path="/divcard_5"><Divcard_5 /></Route>
				<Route exact path="/frame_4"><Frame_4 /></Route>
				<Route exact path="/arrowright"><Arrowright /></Route>
				<Route exact path="/product_card"><Product_card /></Route>
			</Switch>
		</Router>
	);
}
export default RouterDOM;