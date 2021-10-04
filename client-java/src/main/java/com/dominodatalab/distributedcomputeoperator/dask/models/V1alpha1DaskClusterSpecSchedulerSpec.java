/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.dask.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.dask.models.V1alpha1DaskClusterSpecSchedulerAffinityPodAffinityPodAffinityTermLabelSelector;
import com.dominodatalab.distributedcomputeoperator.dask.models.V1alpha1DaskClusterSpecSchedulerResources;
import com.dominodatalab.distributedcomputeoperator.dask.models.V1alpha1DaskClusterSpecSchedulerSpecDataSource;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * V1alpha1DaskClusterSpecSchedulerSpec
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:36:51.907Z[Etc/UTC]")
public class V1alpha1DaskClusterSpecSchedulerSpec {
  public static final String SERIALIZED_NAME_ACCESS_MODES = "accessModes";
  @SerializedName(SERIALIZED_NAME_ACCESS_MODES)
  private List<String> accessModes = null;

  public static final String SERIALIZED_NAME_DATA_SOURCE = "dataSource";
  @SerializedName(SERIALIZED_NAME_DATA_SOURCE)
  private V1alpha1DaskClusterSpecSchedulerSpecDataSource dataSource;

  public static final String SERIALIZED_NAME_RESOURCES = "resources";
  @SerializedName(SERIALIZED_NAME_RESOURCES)
  private V1alpha1DaskClusterSpecSchedulerResources resources;

  public static final String SERIALIZED_NAME_SELECTOR = "selector";
  @SerializedName(SERIALIZED_NAME_SELECTOR)
  private V1alpha1DaskClusterSpecSchedulerAffinityPodAffinityPodAffinityTermLabelSelector selector;

  public static final String SERIALIZED_NAME_STORAGE_CLASS_NAME = "storageClassName";
  @SerializedName(SERIALIZED_NAME_STORAGE_CLASS_NAME)
  private String storageClassName;

  public static final String SERIALIZED_NAME_VOLUME_MODE = "volumeMode";
  @SerializedName(SERIALIZED_NAME_VOLUME_MODE)
  private String volumeMode;

  public static final String SERIALIZED_NAME_VOLUME_NAME = "volumeName";
  @SerializedName(SERIALIZED_NAME_VOLUME_NAME)
  private String volumeName;


  public V1alpha1DaskClusterSpecSchedulerSpec accessModes(List<String> accessModes) {
    
    this.accessModes = accessModes;
    return this;
  }

  public V1alpha1DaskClusterSpecSchedulerSpec addAccessModesItem(String accessModesItem) {
    if (this.accessModes == null) {
      this.accessModes = new ArrayList<>();
    }
    this.accessModes.add(accessModesItem);
    return this;
  }

   /**
   * Get accessModes
   * @return accessModes
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<String> getAccessModes() {
    return accessModes;
  }


  public void setAccessModes(List<String> accessModes) {
    this.accessModes = accessModes;
  }


  public V1alpha1DaskClusterSpecSchedulerSpec dataSource(V1alpha1DaskClusterSpecSchedulerSpecDataSource dataSource) {
    
    this.dataSource = dataSource;
    return this;
  }

   /**
   * Get dataSource
   * @return dataSource
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1DaskClusterSpecSchedulerSpecDataSource getDataSource() {
    return dataSource;
  }


  public void setDataSource(V1alpha1DaskClusterSpecSchedulerSpecDataSource dataSource) {
    this.dataSource = dataSource;
  }


  public V1alpha1DaskClusterSpecSchedulerSpec resources(V1alpha1DaskClusterSpecSchedulerResources resources) {
    
    this.resources = resources;
    return this;
  }

   /**
   * Get resources
   * @return resources
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1DaskClusterSpecSchedulerResources getResources() {
    return resources;
  }


  public void setResources(V1alpha1DaskClusterSpecSchedulerResources resources) {
    this.resources = resources;
  }


  public V1alpha1DaskClusterSpecSchedulerSpec selector(V1alpha1DaskClusterSpecSchedulerAffinityPodAffinityPodAffinityTermLabelSelector selector) {
    
    this.selector = selector;
    return this;
  }

   /**
   * Get selector
   * @return selector
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1DaskClusterSpecSchedulerAffinityPodAffinityPodAffinityTermLabelSelector getSelector() {
    return selector;
  }


  public void setSelector(V1alpha1DaskClusterSpecSchedulerAffinityPodAffinityPodAffinityTermLabelSelector selector) {
    this.selector = selector;
  }


  public V1alpha1DaskClusterSpecSchedulerSpec storageClassName(String storageClassName) {
    
    this.storageClassName = storageClassName;
    return this;
  }

   /**
   * Get storageClassName
   * @return storageClassName
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getStorageClassName() {
    return storageClassName;
  }


  public void setStorageClassName(String storageClassName) {
    this.storageClassName = storageClassName;
  }


  public V1alpha1DaskClusterSpecSchedulerSpec volumeMode(String volumeMode) {
    
    this.volumeMode = volumeMode;
    return this;
  }

   /**
   * Get volumeMode
   * @return volumeMode
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getVolumeMode() {
    return volumeMode;
  }


  public void setVolumeMode(String volumeMode) {
    this.volumeMode = volumeMode;
  }


  public V1alpha1DaskClusterSpecSchedulerSpec volumeName(String volumeName) {
    
    this.volumeName = volumeName;
    return this;
  }

   /**
   * Get volumeName
   * @return volumeName
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getVolumeName() {
    return volumeName;
  }


  public void setVolumeName(String volumeName) {
    this.volumeName = volumeName;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1DaskClusterSpecSchedulerSpec v1alpha1DaskClusterSpecSchedulerSpec = (V1alpha1DaskClusterSpecSchedulerSpec) o;
    return Objects.equals(this.accessModes, v1alpha1DaskClusterSpecSchedulerSpec.accessModes) &&
        Objects.equals(this.dataSource, v1alpha1DaskClusterSpecSchedulerSpec.dataSource) &&
        Objects.equals(this.resources, v1alpha1DaskClusterSpecSchedulerSpec.resources) &&
        Objects.equals(this.selector, v1alpha1DaskClusterSpecSchedulerSpec.selector) &&
        Objects.equals(this.storageClassName, v1alpha1DaskClusterSpecSchedulerSpec.storageClassName) &&
        Objects.equals(this.volumeMode, v1alpha1DaskClusterSpecSchedulerSpec.volumeMode) &&
        Objects.equals(this.volumeName, v1alpha1DaskClusterSpecSchedulerSpec.volumeName);
  }

  @Override
  public int hashCode() {
    return Objects.hash(accessModes, dataSource, resources, selector, storageClassName, volumeMode, volumeName);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1DaskClusterSpecSchedulerSpec {\n");
    sb.append("    accessModes: ").append(toIndentedString(accessModes)).append("\n");
    sb.append("    dataSource: ").append(toIndentedString(dataSource)).append("\n");
    sb.append("    resources: ").append(toIndentedString(resources)).append("\n");
    sb.append("    selector: ").append(toIndentedString(selector)).append("\n");
    sb.append("    storageClassName: ").append(toIndentedString(storageClassName)).append("\n");
    sb.append("    volumeMode: ").append(toIndentedString(volumeMode)).append("\n");
    sb.append("    volumeName: ").append(toIndentedString(volumeName)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}
