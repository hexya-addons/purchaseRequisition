package purchase_requisition

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
)

//import odoo.addons.decimal_precision as dp
func init() {
	h.PurchaseRequisitionType().DeclareModel()

	h.PurchaseRequisitionType().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:    "Agreement Type",
			Required:  true,
			Translate: true,
		},
		"Sequence": models.IntegerField{
			Default: models.DefaultValue(1),
		},
		"Exclusive": models.SelectionField{
			Selection: types.Selection{
				"exclusive": "Select only one RFQ (exclusive)",
				"multiple":  "Select multiple RFQ",
			},
			String:   "Agreement Selection Type",
			Required: true,
			Default:  models.DefaultValue("multiple"),
			Help: "Select only one RFQ (exclusive):  when a purchase order" +
				"is confirmed, cancel the remaining purchase order." +
				"" +
				"                    Select multiple RFQ: allows multiple" +
				"purchase orders. On confirmation of a purchase order it" +
				"does not cancel the remaining orders",
		},
		"QuantityCopy": models.SelectionField{
			Selection: types.Selection{
				"copy": "Use quantities of agreement",
				"none": "Set quantities manually",
			},
			String:   "Quantities",
			Required: true,
			Default:  models.DefaultValue("none"),
		},
		"LineCopy": models.SelectionField{
			Selection: types.Selection{
				"copy": "Use lines of agreement",
				"none": "Do not create RfQ lines automatically",
			},
			String:   "Lines",
			Required: true,
			Default:  models.DefaultValue("copy"),
		},
	})
	h.PurchaseRequisition().DeclareModel()

	h.PurchaseRequisition().Methods().GetPickingIn().DeclareMethod(
		`GetPickingIn`,
		func(rs m.PurchaseRequisitionSet) {
			//        pick_in = self.env.ref('stock.picking_type_in')
			//        if not pick_in:
			//            company = self.env['res.company']._company_default_get(
			//                'purchase.requisition')
			//            pick_in = self.env['stock.picking.type'].search(
			//                [('warehouse_id.company_id', '=', company.id),
			//                 ('code', '=', 'incoming')],
			//                limit=1)
			//        return pick_in
		})
	h.PurchaseRequisition().Methods().GetTypeId().DeclareMethod(
		`GetTypeId`,
		func(rs m.PurchaseRequisitionSet) {
			//        return self.env['purchase.requisition.type'].search([], limit=1)
		})
	h.PurchaseRequisition().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:   "Agreement Reference",
			Required: true,
			NoCopy:   true,
			Default:  func(env models.Environment) interface{} { return env["ir.sequence"].next_by_code() },
		},
		"Origin": models.CharField{
			String: "Source Document",
		},
		"OrderCount": models.IntegerField{
			Compute: h.PurchaseRequisition().Methods().ComputeOrdersNumber(),
			String:  "Number of Orders",
		},
		"VendorId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Vendor",
		},
		"TypeId": models.Many2OneField{
			RelationModel: h.PurchaseRequisitionType(),
			String:        "Agreement Type",
			Required:      true,
			Default:       models.DefaultValue(_get_type_id),
		},
		"OrderingDate": models.DateField{
			String: "Ordering Date",
		},
		"DateEnd": models.DateTimeField{
			String: "Agreement Deadline",
		},
		"ScheduleDate": models.DateField{
			String: "Delivery Date",
			Index:  true,
			Help: "The expected and scheduled delivery date where all the" +
				"products are received",
		},
		"UserId": models.Many2OneField{
			RelationModel: h.User(),
			String:        "Responsible",
			Default:       func(env models.Environment) interface{} { return env.Uid() },
		},
		"Description": models.TextField{},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Required:      true,
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
		},
		"PurchaseIds": models.One2ManyField{
			RelationModel: h.PurchaseOrder(),
			ReverseFK:     "",
			String:        "Purchase Orders",
			//states={'done': [('readonly', True)]}
		},
		"LineIds": models.One2ManyField{
			RelationModel: h.PurchaseRequisitionLine(),
			ReverseFK:     "",
			String:        "Products to Purchase",
			//states={'done': [('readonly', True)]}
			NoCopy: false,
		},
		"ProcurementId": models.Many2OneField{
			RelationModel: h.ProcurementOrder(),
			String:        "Procurement",
			OnDelete:      `set null`,
			NoCopy:        true,
		},
		"WarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Warehouse",
		},
		"State": models.SelectionField{
			Selection: types.Selection{
				"draft":       "Draft",
				"in_progress": "Confirmed",
				"open":        "Bid Selection",
				"done":        "Done",
				"cancel":      "Cancelled",
			},
			String: "Status",
			//track_visibility='onchange'
			Required: true,
			NoCopy:   true,
			Default:  models.DefaultValue("draft"),
		},
		"AccountAnalyticId": models.Many2OneField{
			RelationModel: h.AccountAnalyticAccount(),
			String:        "Analytic Account",
		},
		"PickingTypeId": models.Many2OneField{
			RelationModel: h.StockPickingType(),
			String:        "Picking Type",
			Required:      true,
			Default:       models.DefaultValue(_get_picking_in),
		},
	})
	h.PurchaseRequisition().Methods().ComputeOrdersNumber().DeclareMethod(
		`ComputeOrdersNumber`,
		func(rs h.PurchaseRequisitionSet) h.PurchaseRequisitionData {
			//        for requisition in self:
			//            requisition.order_count = len(requisition.purchase_ids)
		})
	h.PurchaseRequisition().Methods().ActionCancel().DeclareMethod(
		`ActionCancel`,
		func(rs m.PurchaseRequisitionSet) {
			//        for requisition in self:
			//            requisition.purchase_ids.button_cancel()
			//            for po in requisition.purchase_ids:
			//                po.message_post(
			//                    body=_('Cancelled by the agreement associated to this quotation.'))
			//        self.write({'state': 'cancel'})
		})
	h.PurchaseRequisition().Methods().ActionInProgress().DeclareMethod(
		`ActionInProgress`,
		func(rs m.PurchaseRequisitionSet) {
			//        if not all(obj.line_ids for obj in self):
			//            raise UserError(
			//                _('You cannot confirm call because there is no product line.'))
			//        self.write({'state': 'in_progress'})
		})
	h.PurchaseRequisition().Methods().ActionOpen().DeclareMethod(
		`ActionOpen`,
		func(rs m.PurchaseRequisitionSet) {
			//        self.write({'state': 'open'})
		})
	h.PurchaseRequisition().Methods().ActionDraft().DeclareMethod(
		`ActionDraft`,
		func(rs m.PurchaseRequisitionSet) {
			//        self.write({'state': 'draft'})
		})
	h.PurchaseRequisition().Methods().ActionDone().DeclareMethod(
		`
        Generate all purchase order based on selected lines,
should only be called on one agreement at a time
        `,
		func(rs m.PurchaseRequisitionSet) {
			//        if any(purchase_order.state in ['draft', 'sent', 'to approve'] for purchase_order in self.mapped('purchase_ids')):
			//            raise UserError(
			//                _('You have to cancel or validate every RfQ before closing the purchase requisition.'))
			//        self.write({'state': 'done'})
		})
	h.PurchaseRequisitionLine().DeclareModel()

	h.PurchaseRequisitionLine().AddFields(map[string]models.FieldDefinition{
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Filter:        q.PurchaseOk().Equals(True),
			Required:      true,
		},
		"ProductUomId": models.Many2OneField{
			RelationModel: h.ProductUom(),
			String:        "Product Unit of Measure",
		},
		"ProductQty": models.FloatField{
			String: "Quantity",
			//digits=dp.get_precision('Product Unit of Measure')
		},
		"PriceUnit": models.FloatField{
			String: "Unit Price",
			//digits=dp.get_precision('Product Price')
		},
		"QtyOrdered": models.FloatField{
			Compute: h.PurchaseRequisitionLine().Methods().ComputeOrderedQty(),
			String:  "Ordered Quantities",
		},
		"RequisitionId": models.Many2OneField{
			RelationModel: h.PurchaseRequisition(),
			String:        "Purchase Agreement",
			OnDelete:      `cascade`,
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			Related:       `RequisitionId.CompanyId`,
			String:        "Company",
			Stored:        true,
			ReadOnly:      true,
			Default:       func(env models.Environment) interface{} { return env["res.company"]._company_default_get() },
		},
		"AccountAnalyticId": models.Many2OneField{
			RelationModel: h.AccountAnalyticAccount(),
			String:        "Analytic Account",
		},
		"ScheduleDate": models.DateField{
			String: "Scheduled Date",
		},
	})
	h.PurchaseRequisitionLine().Methods().ComputeOrderedQty().DeclareMethod(
		`ComputeOrderedQty`,
		func(rs h.PurchaseRequisitionLineSet) h.PurchaseRequisitionLineData {
			//        for line in self:
			//            total = 0.0
			//            for po in line.requisition_id.purchase_ids.filtered(lambda purchase_order: purchase_order.state in ['purchase', 'done']):
			//                for po_line in po.order_line.filtered(lambda order_line: order_line.product_id == line.product_id):
			//                    if po_line.product_uom != line.product_uom_id:
			//                        total += po_line.product_uom._compute_quantity(
			//                            po_line.product_qty, line.product_uom_id)
			//                    else:
			//                        total += po_line.product_qty
			//            line.qty_ordered = total
		})
	h.PurchaseRequisitionLine().Methods().OnchangeProductId().DeclareMethod(
		`OnchangeProductId`,
		func(rs m.PurchaseRequisitionLineSet) {
			//        if self.product_id:
			//            self.product_uom_id = self.product_id.uom_id
			//            self.product_qty = 1.0
			//        if not self.account_analytic_id:
			//            self.account_analytic_id = self.requisition_id.account_analytic_id
			//        if not self.schedule_date:
			//            self.schedule_date = self.requisition_id.schedule_date
		})
	h.PurchaseOrder().DeclareModel()

	h.PurchaseOrder().AddFields(map[string]models.FieldDefinition{
		"RequisitionId": models.Many2OneField{
			RelationModel: h.PurchaseRequisition(),
			String:        "Purchase Agreement",
			NoCopy:        true,
		},
	})
	h.PurchaseOrder().Methods().OnchangeRequisitionId().DeclareMethod(
		`OnchangeRequisitionId`,
		func(rs m.PurchaseOrderSet) {
			//        if not self.requisition_id:
			//            return
			//        requisition = self.requisition_id
			//        if self.partner_id:
			//            partner = self.partner_id
			//        else:
			//            partner = requisition.vendor_id
			//        payment_term = partner.property_supplier_payment_term_id
			//        currency = partner.property_purchase_currency_id or requisition.company_id.currency_id
			//        FiscalPosition = self.env['account.fiscal.position']
			//        fpos = FiscalPosition.get_fiscal_position(partner.id)
			//        fpos = FiscalPosition.browse(fpos)
			//        self.partner_id = partner.id
			//        self.fiscal_position_id = fpos.id
			//        self.payment_term_id = payment_term.id
			//        self.company_id = requisition.company_id.id
			//        self.currency_id = currency.id
			//        self.origin = requisition.name
			//        self.partner_ref = requisition.name
			//        self.notes = requisition.description
			//        self.date_order = requisition.date_end or fields.Datetime.now()
			//        self.picking_type_id = requisition.picking_type_id.id
			//        if requisition.type_id.line_copy != 'copy':
			//            return
			//        order_lines = []
			//        for line in requisition.line_ids:
			//            # Compute name
			//            product_lang = line.product_id.with_context({
			//                'lang': partner.lang,
			//                'partner_id': partner.id,
			//            })
			//            name = product_lang.display_name
			//            if product_lang.description_purchase:
			//                name += '\n' + product_lang.description_purchase
			//
			//            # Compute taxes
			//            if fpos:
			//                taxes_ids = fpos.map_tax(line.product_id.supplier_taxes_id.filtered(
			//                    lambda tax: tax.company_id == requisition.company_id)).ids
			//            else:
			//                taxes_ids = line.product_id.supplier_taxes_id.filtered(
			//                    lambda tax: tax.company_id == requisition.company_id).ids
			//
			//            # Compute quantity and price_unit
			//            if line.product_uom_id != line.product_id.uom_po_id:
			//                product_qty = line.product_uom_id._compute_quantity(
			//                    line.product_qty, line.product_id.uom_po_id)
			//                price_unit = line.product_uom_id._compute_price(
			//                    line.price_unit, line.product_id.uom_po_id)
			//            else:
			//                product_qty = line.product_qty
			//                price_unit = line.price_unit
			//
			//            if requisition.type_id.quantity_copy != 'copy':
			//                product_qty = 0
			//
			//            # Compute price_unit in appropriate currency
			//            if requisition.company_id.currency_id != currency:
			//                price_unit = requisition.company_id.currency_id.compute(
			//                    price_unit, currency)
			//
			//            # Create PO line
			//            order_lines.append((0, 0, {
			//                'name': name,
			//                'product_id': line.product_id.id,
			//                'product_uom': line.product_id.uom_po_id.id,
			//                'product_qty': product_qty,
			//                'price_unit': price_unit,
			//                'taxes_id': [(6, 0, taxes_ids)],
			//                'date_planned': requisition.schedule_date or fields.Date.today(),
			//                'procurement_ids': [(6, 0, [requisition.procurement_id.id])] if requisition.procurement_id else False,
			//                'account_analytic_id': line.account_analytic_id.id,
			//            }))
			//        self.order_line = order_lines
		})
	h.PurchaseOrder().Methods().ButtonConfirm().DeclareMethod(
		`ButtonConfirm`,
		func(rs m.PurchaseOrderSet) {
			//        res = super(PurchaseOrder, self).button_confirm()
			//        for po in self:
			//            if po.requisition_id.type_id.exclusive == 'exclusive':
			//                others_po = po.requisition_id.mapped(
			//                    'purchase_ids').filtered(lambda r: r.id != po.id)
			//                others_po.button_cancel()
			//                po.requisition_id.action_done()
			//
			//            for element in po.order_line:
			//                if element.product_id == po.requisition_id.procurement_id.product_id:
			//                    element.move_ids.write({
			//                        'procurement_id': po.requisition_id.procurement_id.id,
			//                        'move_dest_id': po.requisition_id.procurement_id.move_dest_id.id,
			//                    })
			//        return res
		})
	h.PurchaseOrder().Methods().Create().Extend(
		`Create`,
		func(rs m.PurchaseOrderSet, vals models.RecordData) {
			//        purchase = super(PurchaseOrder, self).create(vals)
			//        if purchase.requisition_id:
			//            purchase.message_post_with_view('mail.message_origin_link',
			//                                            values={
			//                                                'self': purchase, 'origin': purchase.requisition_id},
			//                                            subtype_id=self.env['ir.model.data'].xmlid_to_res_id('mail.mt_note'))
			//        return purchase
		})
	h.PurchaseOrder().Methods().Write().Extend(
		`Write`,
		func(rs m.PurchaseOrderSet, vals models.RecordData) {
			//        result = super(PurchaseOrder, self).write(vals)
			//        if vals.get('requisition_id'):
			//            self.message_post_with_view('mail.message_origin_link',
			//                                        values={
			//                                            'self': self, 'origin': self.requisition_id, 'edit': True},
			//                                        subtype_id=self.env['ir.model.data'].xmlid_to_res_id('mail.mt_note'))
			//        return result
		})
	h.PurchaseOrderLine().DeclareModel()

	h.PurchaseOrderLine().Methods().OnchangeQuantity().DeclareMethod(
		`OnchangeQuantity`,
		func(rs m.PurchaseOrderLineSet) {
			//        res = super(PurchaseOrderLine, self)._onchange_quantity()
			//        if self.order_id.requisition_id:
			//            for line in self.order_id.requisition_id.line_ids:
			//                if line.product_id == self.product_id:
			//                    if line.product_uom_id != self.product_uom:
			//                        self.price_unit = line.product_uom_id._compute_price(
			//                            line.price_unit, self.product_uom)
			//                    else:
			//                        self.price_unit = line.price_unit
			//                    break
			//        return res
		})
	h.ProductTemplate().DeclareModel()

	h.ProductTemplate().AddFields(map[string]models.FieldDefinition{
		"PurchaseRequisition": models.SelectionField{
			Selection: types.Selection{
				"rfq":     "Create a draft purchase order",
				"tenders": "Propose a call for tenders",
			},
			String:  "Procurement",
			Default: models.DefaultValue("rfq"),
		},
	})
	h.ProcurementOrder().DeclareModel()

	h.ProcurementOrder().AddFields(map[string]models.FieldDefinition{
		"RequisitionId": models.Many2OneField{
			RelationModel: h.PurchaseRequisition(),
			String:        "Latest Requisition",
		},
	})
	h.ProcurementOrder().Methods().MakePo().DeclareMethod(
		`MakePo`,
		func(rs m.ProcurementOrderSet) {
			//        Requisition = self.env['purchase.requisition']
			//        procurements = self.env['procurement.order']
			//        Warehouse = self.env['stock.warehouse']
			//        res = []
			//        for procurement in self:
			//            if procurement.product_id.purchase_requisition == 'tenders':
			//                warehouse_id = Warehouse.search(
			//                    [('company_id', '=', procurement.company_id.id)], limit=1).id
			//                requisition_id = Requisition.create({
			//                    'origin': procurement.origin,
			//                    'date_end': procurement.date_planned,
			//                    'warehouse_id': warehouse_id,
			//                    'company_id': procurement.company_id.id,
			//                    'procurement_id': procurement.id,
			//                    'picking_type_id': procurement.rule_id.picking_type_id.id,
			//                    'line_ids': [(0, 0, {
			//                        'product_id': procurement.product_id.id,
			//                        'product_uom_id': procurement.product_uom.id,
			//                        'product_qty': procurement.product_qty
			//                    })],
			//                })
			//                procurement.message_post(
			//                    body=_("Purchase Requisition created"))
			//                requisition_id.message_post_with_view('mail.message_origin_link',
			//                                                      values={
			//                                                          'self': requisition_id, 'origin': procurement},
			//                                                      subtype_id=self.env['ir.model.data'].xmlid_to_res_id('mail.mt_note'))
			//                procurement.requisition_id = requisition_id
			//                procurements += procurement
			//                res += [procurement.id]
			//        set_others = self - procurements
			//        if set_others:
			//            res += super(ProcurementOrder, set_others).make_po()
			//        return res
		})
}
