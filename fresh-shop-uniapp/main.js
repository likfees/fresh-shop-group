import App from './App'

// #ifndef VUE3
import Vue from 'vue'
import uView from '@/uni_modules/uview-ui'
import pageWrapper from '@/components/pageWrapper/pageWrapper.vue'
import toast from '@/utils/toast.js'
import {filters} from '@/filter/filters.js'
Vue.config.productionTip = false
App.mpType = 'app'

// 定义全局自定义过滤器
Object.keys(filters).forEach(key => {
	Vue.filter(key, filters[key])
})

try {
	function isPromise(obj) {
		return (
			!!obj &&
			(typeof obj === "object" || typeof obj === "function") &&
			typeof obj.then === "function"
		);
	}

	// 统一 vue2 API Promise 化返回格式与 vue3 保持一致
	uni.addInterceptor({
		returnValue(res) {
			if (!isPromise(res)) {
				return res;
			}
			return new Promise((resolve, reject) => {
				res.then((res) => {
					if (res[0]) {
						reject(res[0]);
					} else {
						resolve(res[1]);
					}
				});
			});
		},
	});
} catch (error) {}

Vue.use(uView);
Vue.use(pageWrapper)

Vue.prototype.$message = toast.message

const app = new Vue({
	...App
})
app.$mount()
// #endif

// #ifdef VUE3
import {
	createSSRApp
} from 'vue'
export function createApp() {
	const app = createSSRApp(App)
	return {
		app
	}
}
// #endif
