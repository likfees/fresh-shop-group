<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
           <el-form-item label="售后状态" prop="status">
            <el-select v-model="searchInfo.status" clearable placeholder="请选择" @clear="()=>{searchInfo.status=undefined}">
              <el-option v-for="(item,key) in returnStatusOptions" :key="key" :label="item.label" :value="item.value" />
            </el-select>
            </el-form-item>
           <el-form-item label="退款状态" prop="refundStatus">
            <el-select v-model="searchInfo.refundStatus" clearable placeholder="请选择" @clear="()=>{searchInfo.refundStatus=undefined}">
              <el-option v-for="(item,key) in ReturnRefundStatusOptions" :key="key" :label="item.label" :value="item.value" />
            </el-select>
            </el-form-item>
        <el-form-item label="售后处理时间">
            
            <el-date-picker v-model="searchInfo.startProcessTime" type="datetime" placeholder="搜索条件（起）"></el-date-picker>
            —
            <el-date-picker v-model="searchInfo.endProcessTime" type="datetime" placeholder="搜索条件（止）"></el-date-picker>

        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
           —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
                <el-button type="primary" @click="onDelete">确定</el-button>
            </div>
            <template #reference>
                <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="deleteVisible = true">删除</el-button>
            </template>
            </el-popover>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="编号" prop="ID" width="90" />
        <el-table-column align="left" label="用户id" prop="userId" width="120" />
        <el-table-column align="left" label="订单Id" prop="orderId" width="120" />
        <el-table-column align="left" label="申请原因" prop="reason" width="120" />
        <el-table-column align="left" label="退款金额" prop="amount" width="120" />
        <el-table-column align="left" label="售后状态" prop="status" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.status,returnStatusOptions) }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="退款状态" prop="refundStatus" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.refundStatus,ReturnRefundStatusOptions) }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="售后说明" prop="reply" width="120" />
         <el-table-column align="left" label="售后处理时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.processTime) }}</template>
         </el-table-column>
        <el-table-column align="left" label="创建日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateOrderReturnFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 20, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="type === 'create' ? '新增操作' : '修改操作'">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="用户id:"  prop="userId" >
          <el-input v-model.number="formData.userId" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="订单Id:"  prop="orderId" >
          <el-input v-model.number="formData.orderId" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="申请原因:"  prop="reason" >
          <el-input v-model="formData.reason" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="退款金额:"  prop="amount" >
          <el-input-number v-model="formData.amount"  style="width:100%" :precision="2" :clearable="true"  />
        </el-form-item>
        <el-form-item label="售后状态:"  prop="status" >
          <el-select v-model="formData.status" placeholder="请选择" style="width:100%" :clearable="true" >
            <el-option v-for="(item,key) in returnStatusOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="退款状态:"  prop="refundStatus" >
          <el-select v-model="formData.refundStatus" placeholder="请选择" style="width:100%" :clearable="true" >
            <el-option v-for="(item,key) in ReturnRefundStatusOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="售后说明:"  prop="reply" >
          <el-input v-model="formData.reply" :clearable="true"  placeholder="请输入" />
        </el-form-item>
        <el-form-item label="售后处理时间:"  prop="processTime" >
          <el-date-picker v-model="formData.processTime" type="date" style="width:100%" placeholder="选择日期" :clearable="true"  />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'OrderReturn'
}
</script>

<script setup>
import {
  createOrderReturn,
  deleteOrderReturn,
  deleteOrderReturnByIds,
  updateOrderReturn,
  findOrderReturn,
  getOrderReturnList
} from '@/api/orderReturn'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

// 自动化生成的字典（可能为空）以及字段
const returnStatusOptions = ref([])
const ReturnRefundStatusOptions = ref([])
const formData = ref({
        userId: 0,
        orderId: 0,
        reason: '',
        amount: 0,
        status: undefined,
        refundStatus: undefined,
        reply: '',
        processTime: new Date(),
        })

// 验证规则
const rule = reactive({
})

const elFormRef = ref()


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getOrderReturnList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
    returnStatusOptions.value = await getDictFunc('returnStatus')
    ReturnRefundStatusOptions.value = await getDictFunc('ReturnRefundStatus')
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteOrderReturnFunc(row)
        })
    }


// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async() => {
      const ids = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          ids.push(item.ID)
        })
      const res = await deleteOrderReturnByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        deleteVisible.value = false
        getTableData()
      }
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateOrderReturnFunc = async(row) => {
    const res = await findOrderReturn({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.reorderReturn
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteOrderReturnFunc = async (row) => {
    const res = await deleteOrderReturn({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        userId: 0,
        orderId: 0,
        reason: '',
        amount: 0,
        status: undefined,
        refundStatus: undefined,
        reply: '',
        processTime: new Date(),
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createOrderReturn(formData.value)
                  break
                case 'update':
                  res = await updateOrderReturn(formData.value)
                  break
                default:
                  res = await createOrderReturn(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '操作成功'
                })
                closeDialog()
                getTableData()
              }
      })
}
</script>

<style>
</style>
